package repository

import (
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type ChallengeRepository interface {
	AddChallenge(c *model.Challenge) error
	UpdateChallenge(c *model.Challenge) error

	ListAllChallenges() ([]*model.Challenge, error)
	ListOpenedChallenges() ([]*model.Challenge, error)
	ListChallengeByIDs(ids []uint32) ([]*model.Challenge, error)

	ListAllTags() ([]*model.Tag, error)
	ListAllAttachments() ([]*model.Attachment, error)
	FindTagsByChallengeID(id uint32) ([]*model.Tag, error)
	FindAttachmentsByChallengeID(id uint32) ([]*model.Attachment, error)
	GetChallengeByID(challengeID uint32) (*model.Challenge, error)
	GetChallengeByName(name string) (*model.Challenge, error)
	GetChallengeByFlag(flag string) (*model.Challenge, error)

	AddChallengeAttachment(a *model.Attachment) error
	AddChallengeTag(t *model.Tag) error
	ListTagsByChallengeIDs(ids []uint32) ([]*model.Tag, error)
	ListAttachmentsByChallengeIDs(ids []uint32) ([]*model.Attachment, error)
	DeleteAttachmentByChallengeId(challengeId uint32) error
	DeleteTagByChallengeId(challengeId uint32) error

	OpenChallengeByID(chalelngeID uint32) error
	CloseChallengeByID(chalelngeID uint32) error
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

func (r *repository) ListOpenedChallenges() ([]*model.Challenge, error) {
	var challenges []*model.Challenge
	if err := r.db.Where("is_open = ?", true).Find(&challenges).Error; err != nil {
		return nil, err
	}
	return challenges, nil
}

func (r *repository) ListChallengeByIDs(ids []uint32) ([]*model.Challenge, error) {
	var challenges []*model.Challenge
	if err := r.db.Where("id IN ?", ids).Find(&challenges).Error; err != nil {
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
func (r *repository) FindTagsByChallengeID(id uint32) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.db.Where("challenge_id = ?", id).Find(&tags).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return tags, nil

}
func (r *repository) FindAttachmentsByChallengeID(id uint32) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	if err := r.db.Where("challenge_id = ?", id).Find(&attachments).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return attachments, nil
}

func (r *repository) GetChallengeByID(challengeID uint32) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("id = ?", challengeID).First(&c).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf(": %w", NotFound("challenge"))
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &c, nil
}

func (r *repository) GetChallengeByName(name string) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("name = ?", name).First(&c).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf(": %w", NotFound("challenge"))
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &c, nil
}

func (r *repository) GetChallengeByFlag(flag string) (*model.Challenge, error) {
	var c model.Challenge
	if err := r.db.Where("flag = ?", flag).First(&c).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *repository) ListTagsByChallengeIDs(ids []uint32) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.db.Order("challenge_id asc").Where("challenge_id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return tags, nil
}

func (r *repository) ListAttachmentsByChallengeIDs(ids []uint32) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	if err := r.db.Order("challenge_id asc").Where("challenge_id IN ?", ids).Find(&attachments).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return attachments, nil
}

func (r *repository) DeleteAttachmentByChallengeId(challengeId uint32) error {
	if err := r.db.Where("challenge_id = ?", challengeId).Delete(&model.Attachment{}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
func (r *repository) DeleteTagByChallengeId(challengeId uint32) error {
	if err := r.db.Where("challenge_id = ?", challengeId).Delete(&model.Tag{}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) OpenChallengeByID(challengeID uint32) error {
	if err := r.db.Model(&model.Challenge{}).Where("id = ?", challengeID).Update("is_open", true).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) CloseChallengeByID(challengeID uint32) error {
	if err := r.db.Model(&model.Challenge{}).Where("id = ?", challengeID).Update("is_open", false).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
