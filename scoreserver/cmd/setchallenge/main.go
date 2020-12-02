package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/theoremoon/kosenctfx/scoreserver/config"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"golang.org/x/xerrors"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Uploader interface {
	Upload(name string, data []byte) (string, error)
}

type Challenge struct {
	Name        string
	Description string
	Flag        string
	Author      string
	IsSurvey    bool
	Host        *string
	Port        *int

	Tags []string
}

func run() error {
	flag.Usage = func() {
		fmt.Printf("Usage:\n  %s [OPTIONS] [CHALLENGES]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	path := flag.String("path", "", "path to directory which contains challenges")
	transfersh := flag.String("transfersh", "", "uploading url of transfer.sh")
	flag.Parse()

	if path == nil || *path == "" {
		flag.Usage()
		return nil
	}

	var uploader Uploader
	if transfersh != nil && *transfersh != "" {
		uploader = &transfershUploader{
			url: *transfersh,
		}
	} else {
		flag.Usage()
		return nil
	}

	conf, err := config.Load()
	if err != nil {
		return err
	}
	db, err := gorm.Open("mysql", conf.Dbdsn)
	if err != nil {
		return err
	}
	defer db.Close()
	db.BlockGlobalUpdate(true)
	repo := repository.New(db)
	repo.Migrate()

	// dig into the directory
	challenges := make([]string, 0)
	filepath.Walk(*path, func(p string, i os.FileInfo, err error) error {
		if i.Name() == "challenge.json" {
			challenges = append(challenges, filepath.Dir(p))
			return filepath.SkipDir
		}
		return nil
	})

	// read challenges
	challengeMap := make(map[string]Challenge)
	for _, c := range challenges {
		jsonData, err := ioutil.ReadFile(filepath.Join(c, "challenge.json"))
		if err != nil {
			return err
		}
		var cdata Challenge
		if err := json.Unmarshal(jsonData, &cdata); err != nil {
			return err
		}
		challengeMap[c] = cdata
	}

	// upload attachments
	attachmentMap := make(map[string][]string)
	for _, c := range challenges {
		attachmentMap[c] = make([]string, 0)
		distfiles := filepath.Join(c, "distfiles")
		if st, err := os.Stat(distfiles); err == nil && st.IsDir() {
			func() {
				buf, err := compressDirectory(distfiles)
				if err != nil {
					log.Println(err)
					return
				}

				filename := fmt.Sprintf("%s_%s.tar.gz", challengeMap[c].Name, hexDigest(buf))
				url, err := uploader.Upload(filename, buf)
				if err != nil {
					log.Println(err)
					return
				}
				attachmentMap[c] = append(attachmentMap[c], url)
				log.Printf("UPLOAD %s as %s\n", filename, url)
			}()
		}

		distarchive := filepath.Join(c, "distarchive")
		if st, err := os.Stat(distarchive); err == nil && st.IsDir() {
			filepath.Walk(distarchive, func(p string, i os.FileInfo, err error) error {
				if !i.Mode().IsRegular() {
					return nil
				}
				buf, err := ioutil.ReadFile(p)
				if err != nil {
					return err
				}

				base, ext := splitExt(filepath.Base(p))
				filename := fmt.Sprintf("%s_%s%s", base, hexDigest(buf), ext)
				url, err := uploader.Upload(filename, buf)
				if err != nil {
					log.Println(err)
					return nil
				}
				attachmentMap[c] = append(attachmentMap[c], url)
				log.Printf("UPLOAD %s as %s\n", filename, url)
				return nil
			})
		}
	}

	// register / update challenges
	for _, c := range challenges {
		oldChal, err := repo.FindChallengeByName(challengeMap[c].Name)
		if err == nil {
			// if challenge is already registered
			// 1. remove current attachments and tags
			// 2. re-set attachments and tags
			// 3. update challenge
			if err := repo.DeleteAttachmentByChallengeId(oldChal.ID); err != nil {
				log.Println(err)
				continue
			}
			for _, a := range attachmentMap[c] {
				if err := repo.AddChallengeAttachment(&model.Attachment{
					ChallengeId: oldChal.ID,
					URL:         a,
				}); err != nil {
					log.Println(err)
					continue
				}
			}
			if err := repo.DeleteTagByChallengeId(oldChal.ID); err != nil {
				log.Println(err)
				continue
			}
			for _, t := range challengeMap[c].Tags {
				if err := repo.AddChallengeTag(&model.Tag{
					ChallengeId: oldChal.ID,
					Tag:         t,
				}); err != nil {
					log.Println(err)
					continue
				}
			}
			oldChal.Flag = challengeMap[c].Flag
			oldChal.Description = challengeMap[c].Description
			oldChal.Author = challengeMap[c].Author
			oldChal.Host = challengeMap[c].Host
			oldChal.Port = challengeMap[c].Port
			oldChal.IsSurvey = challengeMap[c].IsSurvey
			repo.UpdateChallenge(oldChal)
		} else if xerrors.Is(err, gorm.ErrRecordNotFound) {
			// if challenge has not registered yet
			// 1. register challenge
			// 2. set atachments and tags
			chal := model.Challenge{
				Name:        challengeMap[c].Name,
				Flag:        challengeMap[c].Flag,
				Description: challengeMap[c].Description,
				Author:      challengeMap[c].Author,
				Host:        challengeMap[c].Host,
				Port:        challengeMap[c].Port,
				IsSurvey:    challengeMap[c].IsSurvey,
				IsOpen:      false,
			}
			if err := repo.AddChallenge(&chal); err != nil {
				log.Println(err)
				continue
			}
			for _, a := range attachmentMap[c] {
				if err := repo.AddChallengeAttachment(&model.Attachment{
					ChallengeId: chal.ID,
					URL:         a,
				}); err != nil {
					log.Println(err)
					continue
				}
			}
			for _, t := range challengeMap[c].Tags {
				if err := repo.AddChallengeTag(&model.Tag{
					ChallengeId: chal.ID,
					Tag:         t,
				}); err != nil {
					log.Println(err)
					continue
				}
			}

		} else {
			log.Printf("%v\n", err)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func compressDirectory(dir string) ([]byte, error) {
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)
	tw := tar.NewWriter(gz)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}
		header.Name, err = filepath.Rel(dir, path)
		if err != nil {
			panic(err)
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})

	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func hexDigest(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func splitExt(path string) (string, string) {
	p := strings.Index(path, ".")
	if p == -1 {
		return path, ""
	}
	return path[:p], path[p:]
}
