package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
)

type transfershUploader struct {
	url string
}

func (uploader *transfershUploader) Upload(name string, data []byte) (string, error) {
	u, err := url.Parse(uploader.url)
	if err != nil {
		return "", err
	}
	u.Path = filepath.Join(u.Path, name)

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return "", errors.New(string(content))
	}

	return string(content), nil
}
