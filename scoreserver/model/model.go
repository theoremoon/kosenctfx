package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint32         `gorm:"primary_key" json:"id"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (model *Model) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New().ID()
	return nil
}

type Team struct {
	Model

	Teamname     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash string
	CountryCode  string

	IsAdmin bool
}

type LoginToken struct {
	Model

	TeamId    uint32
	Token     string `gorm:"unique"`
	IPAddress string
	ExpiresAt int64
}

type PasswordResetToken struct {
	Model

	TeamId    uint32
	Token     string `gorm:"unique"`
	ExpiresAt int64
}

type Challenge struct {
	Model `json:"model"`

	Name        string  `gorm:"unique" json:"name"`
	Flag        string  `gorm:"unique" json:"flag"`
	Description string  `gorm:"size:10000" json:"description"`
	Author      string  `json:"author"`
	Host        *string `json:"host"`
	Port        *int    `json:"port"`

	IsOpen    bool `json:"is_open"`
	IsRunning bool `json:"is_running"`
	IsSurvey  bool `json:"is_survey"`
}

type Tag struct {
	Model

	ChallengeId uint32
	Tag         string
}

type Attachment struct {
	Model

	ChallengeId uint32
	Name        string
	URL         string
}

type Submission struct {
	Model

	ChallengeId *uint32
	TeamId      uint32
	IsCorrect   bool
	Flag        string
	IPAddress   string
}

type ValidSubmission struct {
	Model

	ChallengeId  uint32 `gorm:"unique_index:valid_submission"`
	TeamId       uint32 `gorm:"unique_index:valid_submission"`
	SubmissionId uint32
}

type SubmissionLock struct {
	Model

	TeamId uint32
	Until  int64
}

type Config struct {
	Model

	Token string

	CTFName string
	StartAt int64
	EndAt   int64

	RegisterOpen bool
	CTFOpen      bool

	LockCount    int64
	LockSecond   int
	LockDuration int

	ScoreExpr string `gorm:"size:10000"`
}
