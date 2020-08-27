package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
)

type ChallengeRepository interface {
	AddChallenge(c *model.Challenge) error
	UpdateChallenge(c *model.Challenge) error

	ListAllChallenges() ([]*model.Challenge, error)
	ListAllTags() ([]*model.Tag, error)
	ListAllAttachments() ([]*model.Attachment, error)
	FindTagsByChallengeID(id uint) ([]*model.Tag, error)
	FindAttachmentsByChallengeID(id uint) ([]*model.Attachment, error)
	GetChallengeByID(challengeID uint) (*model.Challenge, error)
	GetChallengeByName(name string) (*model.Challenge, error)
	GetChallengeByFlag(flag string) (*model.Challenge, error)

	AddChallengeAttachment(a *model.Attachment) error
	AddChallengeTag(t *model.Tag) error
	DeleteAttachmentByChallengeId(challengeId uint) error
	DeleteTagByChallengeId(challengeId uint) error

	OpenChallengeByID(chalelngeID uint) error
	CloseChallengeByID(chalelngeID uint) error
}

func (r *repository) AddChallenge(c *model.Challenge) error {
	if err := r.db.Create(c).Error; err != nil {
		if isDuplicatedError(err) {
			return xerrors.Errorf(": %w", Duplicated("challenge"))
		}
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) UpdateChallenge(c *model.Challenge) error {
	if err := r.db.Save(c).Error; err != nil {
		return xerrors.Errorf(": %w", err)
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
func (r *repository) FindTagsByChallengeID(id uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.db.Where("challenge_id = ?", id).Find(&tags).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return tags, nil

}
func (r *repository) FindAttachmentsByChallengeID(id uint) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	if err := r.db.Where("challenge_id = ?", id).Find(&attachments).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return attachments, nil
}

func (r *repository) GetChallengeByID(challengeID uint) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("id = ?", challengeID).First(&c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, xerrors.Errorf(": %w", NotFound("challenge"))
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &c, nil
}

func (r *repository) GetChallengeByName(name string) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("name = ?", name).First(&c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, xerrors.Errorf(": %w", NotFound("challenge"))
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &c, nil
}

func (r *repository) GetChallengeByFlag(flag string) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("flag = ?", flag).First(&c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, xerrors.Errorf(": %w", NotFound("challenge"))
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &c, nil
}

func (r *repository) AddChallengeAttachment(a *model.Attachment) error {
	if err := r.db.Create(a).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
func (r *repository) AddChallengeTag(t *model.Tag) error {
	if err := r.db.Create(t).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) DeleteAttachmentByChallengeId(challengeId uint) error {
	if err := r.db.Where("challenge_id = ?", challengeId).Delete(&model.Attachment{}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
func (r *repository) DeleteTagByChallengeId(challengeId uint) error {
	if err := r.db.Where("challenge_id = ?", challengeId).Delete(&model.Tag{}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) OpenChallengeByID(challengeID uint) error {
	if err := r.db.Model(&model.Challenge{}).Where("id = ?", challengeID).Update("is_open", true).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) CloseChallengeByID(challengeID uint) error {
	if err := r.db.Model(&model.Challenge{}).Where("id = ?", challengeID).Update("is_open", false).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
