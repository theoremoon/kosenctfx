package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
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

type discordWebhook struct {
	sync.Mutex
	url            string
	cancel         context.CancelFunc
	waitDuration   time.Duration
	currentMessage string
}

func NewDiscord(url string, waitDuration time.Duration) Webhook {
	return &discordWebhook{
		url:          url,
		waitDuration: waitDuration,
	}
}

func (w *discordWebhook) Post(message string) error {
	w.Lock()
	if w.cancel != nil {
		w.cancel()
	}
	w.Unlock()

	if w.currentMessage != "" {
		w.currentMessage += "\n" + message
	} else {
		w.currentMessage = message
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	go w.waitAndPost(ctx)

	return nil
}

func (w *discordWebhook) waitAndPost(ctx context.Context) {
	defer func() {
		w.Lock()
		w.cancel = nil
		w.Unlock()
	}()

	select {
	case <-time.After(w.waitDuration):
		if err := w.doPost(); err != nil {
		}
	case <-ctx.Done():
	}
}

func (w *discordWebhook) doPost() error {
	if w.currentMessage == "" {
		return nil
	}

	payload, err := json.Marshal(map[string]interface{}{
		"content": w.currentMessage,
	})
	w.currentMessage = ""
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
