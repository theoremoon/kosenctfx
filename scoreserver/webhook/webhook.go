package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type Webhook interface {
	Post(message string) error
}

type dummyWebhook struct {
	prefix string
}

func Dummy(prefix string) Webhook {
	return &dummyWebhook{
		prefix: prefix,
	}
}

func (w *dummyWebhook) Post(message string) error {
	log.Printf("WEBHOOK:[%s] %s\n", w.prefix, message)
	return nil
}

type webhook struct {
	url string
}

func New(url string) Webhook {
	return &webhook{
		url: url,
	}
}

func (w *webhook) Post(message string) error {
	payload, err := json.Marshal(map[string]interface{}{
		"content": message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", w.url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		data, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(data))
	}

	return nil
}
