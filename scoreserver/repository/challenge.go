package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type ChallengeRepository interface {
	AddChallenge(c *model.Challenge) error
	UpdateChallenge(c *model.Challenge) error

	ListAllChallenges() ([]*model.Challenge, error)
	ListAllTags() ([]*model.Tag, error)
	ListAllAttachments() ([]*model.Attachment, error)
	FindChallengeByName(name string) (*model.Challenge, error)
	FindChallengeByFlag(flag string) (*model.Challenge, error)

	AddChallengeAttachment(a *model.Attachment) error
	AddChallengeTag(t *model.Tag) error
	DeleteAttachmentByChallengeId(challengeId uint) error
	DeleteTagByChallengeId(challengeId uint) error
}

func (r *repository) AddChallenge(c *model.Challenge) error {
	if err := r.db.Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateChallenge(c *model.Challenge) error {
	if err := r.db.Save(c).Error; err != nil {
		return err
	}
	return nil
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

func (r *repository) FindChallengeByName(name string) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("name = ?", name).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) FindChallengeByFlag(flag string) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("flag = ?", flag).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) AddChallengeAttachment(a *model.Attachment) error {
	if err := r.db.Create(a).Error; err != nil {
		return err
	}
	return nil
}
func (r *repository) AddChallengeTag(t *model.Tag) error {
	if err := r.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteAttachmentByChallengeId(challengeId uint) error {
	if err := r.db.Where("challenge_id = ?", challengeId).Delete(&model.Attachment{}).Error; err != nil {
		return err
	}
	return nil
}
func (r *repository) DeleteTagByChallengeId(challengeId uint) error {
	if err := r.db.Where("challenge_id = ?", challengeId).Delete(&model.Tag{}).Error; err != nil {
		return err
	}
	return nil
}
