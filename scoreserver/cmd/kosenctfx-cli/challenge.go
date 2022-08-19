package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/theoremoon/kosenctfx/scoreserver/task"
	"github.com/theoremoon/kosenctfx/scoreserver/task/imagebuilder"
	"gopkg.in/yaml.v2"
)

func challengeMain(dir string) error {
	// dir以下を走査してtask.ymlがあったら読んでいく
	// タスクのあるディレクトリを記録しておく
	taskDefs := make(map[string]*task.TaskDefinition)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "walkdir: ")
		}
		if d.Name() != "task.yml" {
			return nil
		}

		taskdef, err := task.LoadTaskDefinition(path)
		if err != nil {
			return err
		}
		taskDefs[filepath.Dir(path)] = taskdef

		// task.yml があるディレクトリをこれより深く読んでも無駄
		return fs.SkipDir
	})
	if err != nil {
		return err
	}

	// docker imageをbuilderするためのbuilderを準備する
	// imagebuilderにはdocker imageをするための認証情報をサーバからもらってきて渡す
	ctx := context.Background()
	registryConf, err := client.GetRegistryConf(ctx)
	if err != nil {
		return err
	}

	builder, err := imagebuilder.New(&imagebuilder.RegistryConfig{
		URL:      registryConf.URL,
		Username: registryConf.Username,
		Password: registryConf.Password,
	})
	if err != nil {
		return err
	}

	// タスクをどんどんbuild and push and set to scoreserver
	for dir, taskdef := range taskDefs {
		err := uploadTask(ctx, builder, dir, taskdef)
		if err != nil {
			log.Printf("[-] error %s: %v\n", filepath.Base(dir), err)
			continue
		}
		log.Printf("[+] updated %s\n", filepath.Base(dir))
	}

	return nil
}

// 個別のタスクを読みこんで、docker-compose.ymlを読んだり
// 配布するファイルをアップロードしたりする
func uploadTask(ctx context.Context, builder imagebuilder.ImageBuilder, dir string, taskDef *task.TaskDefinition) error {
	id := filepath.Base(dir)
	if taskDef.ID == "" {
		taskDef.ID = id
	}

	attachments, err := prepareAttachments(ctx, id, dir)
	if err != nil {
		return err
	}
	taskDef.Attachments = attachments

	// そのうちcompose.yaml に対応する
	composePath := filepath.Join(dir, "docker-compose.yml")
	if _, err := os.Stat(composePath); !os.IsNotExist(err) {
		data, err := os.ReadFile(composePath)
		if err != nil {
			return err
		}

		config, err := task.ParseComposeConfig(composePath, data)
		if err != nil {
			return err
		}
		err = task.ValidateComposeConfig(config)
		if err != nil {
			return err
		}

		// docker build
		// このタイミングでconfigが書き換わってる（！）
		err = builder.BuildAndPush(ctx, taskDef, config)
		if err != nil {
			return err
		}

		composeBuf, err := yaml.Marshal(config)
		if err != nil {
			return err
		}
		taskDef.Compose = string(composeBuf)
	}

	err = client.NewChallenge(ctx, taskDef)
	if err != nil {
		return err
	}

	return nil
}

// distfiles/ 以下を${id}_${md5}.tar.gz にかためてアップロードする
// rawdistfiles/ 直下のファイルをそのままアップロードする
func prepareAttachments(ctx context.Context, id, dir string) ([]task.Attachment, error) {
	attachments := make([]task.Attachment, 0)

	// distfiles
	distDir := filepath.Join(dir, "distfiles")
	stat, err := os.Stat(distDir)
	if !os.IsNotExist(err) && stat.IsDir() {
		tar, err := makeTar(distDir, id)
		if err != nil {
			return nil, err
		}
		md5sum := md5.Sum(tar)
		filename := fmt.Sprintf("%s_%s.tar.gz", id, hex.EncodeToString(md5sum[:]))
		url, err := uploadFile(ctx, filename, tar)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, task.Attachment{
			URL:  url,
			Name: filename,
		})
	}

	// rawdistfiles
	rawDistdir := filepath.Join(dir, "rawdistfiles")
	stat, err = os.Stat(rawDistdir)
	if !os.IsNotExist(err) && stat.IsDir() {
		filepath.WalkDir(rawDistdir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			url, err := uploadFile(ctx, d.Name(), data)
			if err != nil {
				return err
			}
			attachments = append(attachments, task.Attachment{
				URL:  url,
				Name: d.Name(),
			})
			return nil
		})
	}

	return attachments, nil
}

// distfiles/ 以下のファイルをtar.gzにまとめる
// /bin/sh と tar にパスが通ってることを仮定してる
func makeTar(dir, name string) ([]byte, error) {
	transform := fmt.Sprintf(" --transform 's:^\\./:./%s/:'", name)
	cmd := exec.Command("sh", "-c", "find . -type f | tar cz --files-from=- --to-stdout  --sort=name"+transform)
	cmd.Dir = dir
	tardata, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return tardata, nil
}

func uploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	res, err := client.GetPresignedURL(ctx, filename)
	if err != nil {
		return "", err
	}
	if err := putData(ctx, res.PresignedURL, data); err != nil {
		return "", err
	}
	return res.DownloadURL, nil
}

// PUTメソッドで単にデータを送りつける
// curl -X PUT url --data @file 相当
func putData(ctx context.Context, url string, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "NewRequest")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer res.Body.Close()
	return nil
}
