package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (model *Model) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().ID())
}

type Team struct {
	Model

	Teamname     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash string

	IsAdmin bool
}

type LoginToken struct {
	Model

	TeamId    uint
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
}

type PasswordResetToken struct {
	Model

	TeamId    uint
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
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

	ChallengeId uint
	Tag         string
}

type Attachment struct {
	Model

	ChallengeId uint
	Name        string
	URL         string
}

type Submission struct {
	Model

	ChallengeId *uint
	TeamId      uint
	IsCorrect   bool
	Flag        string
}

type ValidSubmission struct {
	Model

	ChallengeId  uint `gorm:"unique_index:valid_submission"`
	TeamId       uint `gorm:"unique_index:valid_submission"`
	SubmissionId uint
}

type SubmissionLock struct {
	Model

	TeamId uint
	Until  time.Time
}

type Notification struct {
	Model

	Title   string
	Content string
}

type Config struct {
	Model

	Token string

	CTFName string
	StartAt time.Time
	EndAt   time.Time

	RegisterOpen bool
	CTFOpen      bool

	LockCount    int
	LockSecond   int
	LockDuration int

	ScoreExpr string `gorm:"size:10000"`
}
