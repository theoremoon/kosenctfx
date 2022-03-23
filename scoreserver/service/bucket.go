package service

import (
	"github.com/pkg/errors"
	"github.com/theoremoon/kosenctfx/scoreserver/bucket"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
)

type BucketApp interface {
	SetBucketConfig(ctfConf *model.Config, conf *model.BucketConfig) error
	GetBucketClient(ctfConf *model.Config) (bucket.Bucket, error)
	GetBucketConfig(ctfConf *model.Config) (*model.BucketConfig, error)
}

func (app *app) SetBucketConfig(ctfConf *model.Config, conf *model.BucketConfig) error {
	_, err := bucket.SetupBucket(
		conf.Endpoint,
		conf.Region,
		conf.BucketName,
		conf.AccessKey,
		conf.SecretKey,
		conf.HTTPS,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	conf.CTFId = ctfConf.ID
	err = app.repo.SetBucketConfig(conf)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (app *app) GetBucketClient(ctfConf *model.Config) (bucket.Bucket, error) {
	bucketConfig, err := app.repo.GetBucketConfig(ctfConf.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bucket, err := bucket.New(
		bucketConfig.Endpoint,
		bucketConfig.Region,
		bucketConfig.BucketName,
		bucketConfig.AccessKey,
		bucketConfig.SecretKey,
		bucketConfig.HTTPS,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bucket, nil
}

func (app *app) GetBucketConfig(ctfConf *model.Config) (*model.BucketConfig, error) {
	bucketConf, err := app.repo.GetBucketConfig(ctfConf.ID)
	if err != nil {
		if errors.As(err, &repository.NotFoundError{}) {
			return nil, NewErrorMessage("Bucket is not configured yet")
		}
		return nil, errors.WithStack(err)
	}
	return bucketConf, nil
}
