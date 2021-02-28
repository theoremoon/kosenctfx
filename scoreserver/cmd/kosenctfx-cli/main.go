package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/mod/sumdb/dirhash"
	"gopkg.in/yaml.v2"
)

type TaskYaml struct {
	Name        string
	Description string
	Flag        string
	Author      string
	Tags        []string
	Attachments []service.Attachment
	Host        *string
	Port        *int
	IsSurvey    bool `yaml:"is_survey"`
}

func uploadFile(url, token, filename string, blob []byte) (string, error) {
	type Data struct {
		PresignedURL string `json:"presignedURL"`
		DownloadURL  string `json:"downloadURL"`
	}

	var data Data
	_, err := resty.New().SetAuthToken(token).R().
		SetBody(map[string]interface{}{"key": filename}).
		SetResult(&data).
		Post(url + "/admin/get-presigned-url")
	if err != nil {
		return "", err
	}

	_, err = resty.New().R().
		SetBody(blob).
		Put(data.PresignedURL)
	if err != nil {
		return "", err
	}

	return data.DownloadURL, nil
}

func setChallenge(url, token string, taskInfo TaskYaml) error {
	client := resty.New().SetAuthToken(token)
	_, err := client.R().
		SetBody(taskInfo).
		Post(url + "/admin/new-challenge")
	if err != nil {
		return err
	}
	return nil
}

func loadTaskYaml(path string) (*TaskYaml, error) {
	// load task.yml
	taskb, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tasky TaskYaml
	if err := yaml.Unmarshal(taskb, &tasky); err != nil {
		return nil, err
	}

	hostStr := ""
	if tasky.Host != nil {
		hostStr = *tasky.Host
	}
	portStr := ""
	if tasky.Port != nil {
		portStr = strconv.FormatInt(int64(*tasky.Port), 10)
	}

	r := strings.NewReplacer("{host}", hostStr, "{port}", portStr)
	tasky.Description = r.Replace(tasky.Description)
	return &tasky, nil
}

func makeDistfiles(dir, name string) ([]byte, error) {
	transform := fmt.Sprintf(" --transform 's:^\\./:./%s/:'", name)
	cmd := exec.Command("sh", "-c", "find . -type f | tar cz --files-from=- --to-stdout  --sort=name"+transform)
	cmd.Dir = dir
	tardata, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return tardata, nil
}

func run() error {
	var url, token, dir, hashfile string
	flag.StringVar(&url, "url", "", "An endpoint of scoreserver")
	flag.StringVar(&token, "token", "", "An administrative token")
	flag.StringVar(&dir, "dir", "", "tasks directory")
	flag.StringVar(&hashfile, "hashfile", "", "hash file")
	flag.Usage = func() {
		fmt.Printf("Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if url == "" || token == "" || dir == "" {
		flag.Usage()
		return nil
	}
	url = strings.TrimSuffix(url, "/") // remove trailing /

	hash_entries := make(map[string]string)

	if _, err := os.Stat(hashfile); err == nil {
		data, err := ioutil.ReadFile(hashfile)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &hash_entries); err != nil {
			return err
		}
	}

	targets := make(map[string]*TaskYaml)
	// walk tasks directory
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() != "task.yml" {
			return nil
		}
		tasky, err := loadTaskYaml(path)
		if err != nil {
			return err
		}

		// hash tableに乗っていない OR 更新されていたらtargetsに乗せる
		dirpath := filepath.Dir(path)
		h1, _ := dirhash.HashDir(dirpath, "", dirhash.Hash1)
		h2, exist := hash_entries[tasky.Name]
		if !exist || h1 != h2 {
			hash_entries[tasky.Name] = h1
			targets[dirpath] = tasky

		} else {
			log.Printf("[+] SKIP: %s\n", tasky.Name)
		}

		// このディレクトリは深堀りしない
		return filepath.SkipDir
	})
	if err != nil {
		return err
	}

	for d, tasky := range targets {
		taskID := filepath.Base(d)
		attachments := make([]service.Attachment, 0, 10)
		err = func() error {
			distdir := filepath.Join(d, "distfiles")
			if _, err := os.Stat(distdir); err != nil {
				return nil
			}
			tardata, err := makeDistfiles(distdir, taskID)
			md5sum := md5.Sum(tardata)
			filename := fmt.Sprintf("%s_%s.tar.gz", taskID, hex.EncodeToString(md5sum[:]))
			dlUrl, err := uploadFile(url, token, filename, tardata)
			if err != nil {
				return err
			}
			attachments = append(attachments, service.Attachment{
				URL:  dlUrl,
				Name: filename,
			})
			return nil
		}()
		if err != nil {
			return err
		}

		err = func() error {
			rawDistdir := filepath.Join(d, "rawdistfiles")
			if _, err := os.Stat(rawDistdir); err != nil {
				return nil
			}
			err := filepath.Walk(rawDistdir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				blob, err := ioutil.ReadFile(path)
				if err != nil {
					return nil
				}

				dlUrl, err := uploadFile(url, token, info.Name(), blob)
				if err != nil {
					return err
				}
				attachments = append(attachments, service.Attachment{
					URL:  dlUrl,
					Name: info.Name(),
				})
				return nil
			})
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
		tasky.Attachments = attachments
		if err := setChallenge(url, token, *tasky); err != nil {
			return err
		}

		log.Printf("[+] %s\n", tasky.Name)
	}

	// save
	hashb, err := json.Marshal(hash_entries)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(hashfile, hashb, 0755); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
