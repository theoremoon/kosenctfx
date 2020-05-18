package main

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Uploader struct {
	uploader *s3manager.Uploader
	bucket   string
}

func (s *s3Uploader) Upload(name string, data []byte) (string, error) {
	res, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(name),
		Body:   bytes.NewReader(data),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}
	return res.Location, nil
}
