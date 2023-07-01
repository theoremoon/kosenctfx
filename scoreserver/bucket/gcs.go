package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
	"golang.org/x/xerrors"
)

const (
	objectReader = iam.RoleName("roles/storage.legacyObjectReader")
)

type credential struct {
	PrivateKey string `json:"private_key"`
	Email      string `json:"client_email"`
}

type gcsBucket struct {
	bucketName string
	endpoint   string
	region     string
	insecure   bool
	client     *storage.Client
	creds      credential
}

func NewGCSBucket(bucketName, endpoint, region string, insecure bool) (Bucket, error) {
	ctx := context.Background()

	credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	data, err := os.ReadFile(credFile)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	var creds credential
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return &gcsBucket{
		bucketName: bucketName,
		endpoint:   endpoint,
		region:     region,
		insecure:   insecure,
		client:     client,
		creds:      creds,
	}, nil
}

func (b *gcsBucket) CreateBucket() error {
	bucket := b.client.Bucket(b.bucketName)
	_, err := bucket.Objects(context.Background(), nil).Next()
	if err != nil {
		log.Printf("%+v", err)
	}

	iamHandler := bucket.IAM()

	policy, err := iamHandler.Policy(context.Background())
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if !policy.HasRole(iam.AllUsers, objectReader) {
		policy.Add(iam.AllUsers, objectReader)
		if err := iamHandler.SetPolicy(context.Background(), policy); err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	return nil
}

func (b *gcsBucket) GeneratePresignedURL(key string) (string, map[string]string, string, error) {
	postPolicy, err := storage.GenerateSignedPostPolicyV4(b.bucketName, key, &storage.PostPolicyV4Options{
		GoogleAccessID: b.creds.Email,
		PrivateKey:     []byte(b.creds.PrivateKey),
		Expires:        time.Now().Add(PresignURLLifetime),
		Insecure:       b.insecure,
		Hostname:       b.endpoint,
	})
	if err != nil {
		return "", nil, "", xerrors.Errorf(": %w", err)
	}

	return postPolicy.URL, postPolicy.Fields, b.buildKeyURL(key), nil
}

func (b *gcsBucket) buildEndpoint() string {
	if b.insecure {
		return "http://" + b.endpoint
	}
	return "https://" + b.endpoint
}

func (b *gcsBucket) buildKeyURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", b.buildEndpoint(), b.bucketName, key)
}
