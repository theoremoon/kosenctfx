package bucket

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"golang.org/x/xerrors"
)

var (
	PresignURLLifetime = 10 * time.Minute
)

const bucketPolicyTempalte = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AddPerm\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": [\"s3:GetObject\"], \"Resource\": [\"arn:aws:s3:::%s/*\"]}]}"

type Bucket interface {
	CreateBucket() error
	GeneratePresignedURL(key string) (string, string, error)
}

type s3Bucket struct {
	bucketName string
	endpoint   string
	region     string
	insecure   bool
	accessKey  string
	secretKey  string
}

func NewS3Bucket(bucketName, endpoint, region, accessKey, secretKey string, insecure bool) Bucket {
	return &s3Bucket{
		bucketName: bucketName,
		endpoint:   endpoint,
		region:     region,
		accessKey:  accessKey,
		secretKey:  secretKey,
		insecure:   insecure,
	}
}

// TODO CORS
func (b *s3Bucket) CreateBucket() error {
	cred := credentials.NewStaticCredentials(b.accessKey, b.secretKey, "")
	s, err := session.NewSession(&aws.Config{
		Region:           &b.region,
		Endpoint:         b.buildEndpoint(),
		Credentials:      cred,
		DisableSSL:       &b.insecure,
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	svc := s3.New(s)

	// Bucketを作ってみる
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: &b.bucketName,
	})
	if aerr, ok := err.(awserr.Error); ok && (aerr.Code() == s3.ErrCodeBucketAlreadyExists || aerr.Code() == s3.ErrCodeBucketAlreadyOwnedByYou) {
		// すでにbucketがある場合は何もしない
		return nil
	}
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// BucketにPolicyを与える
	if _, err := svc.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: &b.bucketName,
		Policy: aws.String(fmt.Sprintf(bucketPolicyTempalte, b.bucketName)),
	}); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (b *s3Bucket) GeneratePresignedURL(key string) (string, string, error) {
	cred := credentials.NewStaticCredentials(b.accessKey, b.secretKey, "")
	s, err := session.NewSession(&aws.Config{
		Region:           &b.region,
		Endpoint:         b.buildEndpoint(),
		Credentials:      cred,
		DisableSSL:       &b.insecure,
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return "", "", xerrors.Errorf(": %w", err)
	}

	svc := s3.New(s)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &b.bucketName,
		Key:    &key,
	})
	url, err := req.Presign(PresignURLLifetime)
	if err != nil {
		return "", "", xerrors.Errorf(": %w", err)
	}

	return url, b.buildKeyURL(key), nil
}

func (b *s3Bucket) buildEndpoint() *string {
	if b.insecure {
		return aws.String("http://" + b.endpoint)
	}
	return aws.String("https://" + b.endpoint)
}

func (b *s3Bucket) buildKeyURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", *b.buildEndpoint(), b.bucketName, key)
}
