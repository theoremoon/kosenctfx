package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type ChallengeRepository interface {
	ListAllChallenges() ([]*model.Challenge, error)
	ListAllTags() ([]*model.Tag, error)
	ListAllAttachments() ([]*model.Attachment, error)
	FindChallengeByFlag(flag string) (*model.Challenge, error)
}

func (r *repository) ListAllChallenges() ([]*model.Challenge, error) {
	var challenges []*model.Challenge
	if err := r.db.Find(&challenges).Error; err != nil {
		return nil, err
	}
	return challenges, nil
}

func (r *repository) ListAllTags() ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *repository) ListAllAttachments() ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	if err := r.db.Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *repository) FindChallengeByFlag(flag string) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("flag = ?", flag).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
