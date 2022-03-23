package bucket

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

type Bucket interface {
}

type bucket struct {
	Endpoint   string
	Region     string
	BucketName string
	AccessKey  string
	SecretKey  string
	HTTPS      bool
	client     *minio.Client
}

func New(endpoint, region, bucketName, accessKey, secretKey string, https bool) (Bucket, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &bucket{
		Endpoint:   endpoint,
		Region:     region,
		BucketName: bucketName,
		HTTPS:      https,
		client:     client,
	}, nil
}

func SetupBucket(endpoint, region, bucketName, accessKey, secretKey string, https bool) (Bucket, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	exist, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !exist {
		// create bucket
		err = client.MakeBucket(
			context.Background(),
			bucketName,
			minio.MakeBucketOptions{
				Region: region,
			},
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	bucketPolicy := fmt.Sprintf(`{"Version": "2012-10-17", "Statement": [{"Sid": "AddPerm", "Effect": "Allow", "Principal": "*", "Action": ["s3:GetObject"], "Resource": ["arn:aws:s3:::%s/*"]}]}`, bucketName)

	// set bucket policy
	err = client.SetBucketPolicy(
		context.Background(),
		bucketName,
		bucketPolicy,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &bucket{
		Endpoint:   endpoint,
		Region:     region,
		BucketName: bucketName,
		HTTPS:      https,
		client:     client,
	}, nil
}

func (b *bucket) Upload(objectName string, buf *bytes.Buffer) error {
	_, err := b.client.PutObject(
		context.Background(),
		b.BucketName,
		objectName,
		buf,
		int64(buf.Len()),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *bucket) TestBucket() error {
	// test upload & download
	testObjectName := fmt.Sprintf("test-object-%s", uuid.New().String())
	testData := []byte("test")
	err := b.Upload(testObjectName, bytes.NewBuffer(testData))
	if err != nil {
		return errors.WithStack(err)
	}

	url := b.BuildEndpoint(testObjectName)
	data, err := download(url)
	if err != nil {
		return fmt.Errorf("failed to download test object from bucket: %w", err)
	}

	if bytes.Compare(data, testData) != 0 {
		return errors.New("object integrity is not complete. the bucket may be corrupted")
	}
	return nil
}

func (b *bucket) BuildEndpoint(objectName string) string {
	protocol := "https://"
	if !b.HTTPS {
		protocol = "http://"
	}
	return fmt.Sprintf("%s%s/%s/%s", protocol, b.Endpoint, b.BucketName, objectName)
}

func download(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return buf, nil
}
