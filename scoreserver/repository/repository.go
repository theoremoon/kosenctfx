package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type Repository interface {
	UserRepository
	TeamRepository
	ChallengeRepository
	SubmissionRepository
	ConfigRepository
	Migrate()
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Migrate() {
	r.db.AutoMigrate(
		&model.User{},
		&model.LoginToken{},
		&model.PasswordResetToken{},
		&model.Team{},
		&model.Challenge{},
		&model.Tag{},
		&model.Attachment{},
		&model.Submission{},
		&model.Config{},
	)
}
