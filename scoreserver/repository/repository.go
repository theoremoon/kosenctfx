package repository

import (
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
)

type Repository interface {
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
		&model.LoginToken{},
		&model.PasswordResetToken{},
		&model.Team{},
		&model.Challenge{},
		&model.Tag{},
		&model.Attachment{},
		&model.Submission{},
		&model.ValidSubmission{},
		&model.SubmissionLock{},
		&model.Config{},
	)
}

func isDuplicatedError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if xerrors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1062 {
			return true
		}
	}
	return false
}
