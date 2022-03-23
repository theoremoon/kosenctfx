package repository

import (
	"github.com/pkg/errors"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConfigRepository interface {
	SetConfig(conf *model.Config) error
	GetConfig() (*model.Config, error)

	SetBucketConfig(bucketConfig *model.BucketConfig) error
	GetBucketConfig(ctfID uint32) (*model.BucketConfig, error)
}

func (r *repository) SetConfig(conf *model.Config) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var c model.Config
		if err := tx.First(&c).Error; err != nil {
			if !xerrors.Is(err, gorm.ErrRecordNotFound) {
				return xerrors.Errorf(": %w", err)
			}

			// 存在しない時：つくる
			if err := tx.Create(conf).Error; err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		}

		// 存在するとき： update
		conf.Model = c.Model
		if err := tx.Save(conf).Error; err != nil {
			return xerrors.Errorf(": %w", err)
		}
		return nil
	})
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) GetConfig() (*model.Config, error) {
	var conf model.Config
	if err := r.db.First(&conf).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &conf, nil
}

func (r *repository) SetBucketConfig(bucketConfig *model.BucketConfig) error {
	dbErr := r.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(bucketConfig)
	if err := dbErr.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *repository) GetBucketConfig(ctfID uint32) (*model.BucketConfig, error) {
	var b model.BucketConfig
	if err := r.db.Where("ctf_id = ?", ctfID).First(&b).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFound(err.Error())
		}
		return nil, err
	}
	return &b, nil
}
