package bucket

import (
	"time"
)

var (
	PresignURLLifetime = 10 * time.Minute
)

type Bucket interface {
	CreateBucket() error
	GeneratePresignedURL(key string) (string, map[string]string, string, error) // presignedURL, multipart-form, downloadURL, error
}
