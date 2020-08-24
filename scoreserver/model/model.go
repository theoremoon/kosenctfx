package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (model *Model) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().ID())
}

type User struct {
	Model

	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string `gorm:"unique"`
	IsAdmin      bool

	IconPath *string
	TeamId   uint
}

type LoginToken struct {
	Model

	UserId    uint
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
}

type PasswordResetToken struct {
	Model

	UserId    uint
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
}

type Team struct {
	Model

	Teamname string `gorm:"unique"`
	Token    string `gorm:"unique"`

	IsAdmin bool
}

type Challenge struct {
	Model

	Name        string `gorm:"unique"`
	Flag        string `gorm:"unique"`
	Description string `gorm:"size:10000"`
	Author      string
	Host        *string
	Port        *int

	IsOpen    bool
	IsRunning bool
	IsSurvey  bool
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
	UserId      uint
	TeamId      uint
	IsCorrect   bool
	IsValid     bool
	Flag        string
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

	LockCount     int
	LockFrequency int
	LockDuration  int

	ScoreExpr string `gorm:"size:10000"`
}
