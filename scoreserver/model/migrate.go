package model

import (
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&LoginToken{},
		&PasswordResetToken{},
		&Team{},
		&Challenge{},
		&Tag{},
		&Attachment{},
		&Submission{},
		&ValidSubmission{},
		&SubmissionLock{},
		&Message{},
		&Config{},
	)
	if err != nil {
		return xerrors.Errorf("migrate: %w", err)
	}
	return nil
}
