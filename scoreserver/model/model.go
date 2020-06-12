package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string `gorm:"unique"`
	IsAdmin      bool

	IconPath *string
	TeamId   uint
}

type LoginToken struct {
	gorm.Model

	UserId    uint
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
}

type PasswordResetToken struct {
	gorm.Model

	UserId    uint
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
}

type Team struct {
	gorm.Model

	Teamname string `gorm:"unique"`
	Token    string `gorm:"unique"`

	IsAdmin bool
}

type Challenge struct {
	gorm.Model

	Name        string `gorm:"unique"`
	Flag        string `gorm:"unique"`
	Description string `gorm:""`
	Author      string
	Host        *string
	Port        *int

	IsOpen   bool
	IsSurvey bool
}

type Tag struct {
	gorm.Model

	ChallengeId uint
	Tag         string
}

type Attachment struct {
	gorm.Model

	ChallengeId uint
	URL         string
}

type Submission struct {
	gorm.Model

	ChallengeId *uint
	UserId      uint
	TeamId      uint
	IsCorrect   bool
	IsValid     bool
	Flag        string
}

type Notification struct {
	gorm.Model
	Content string
}

type Qualification struct {
	gorm.Model
	Content      string
	UserId       uint
	ResponseType uint
	Completed    bool
	IsPublic     bool
}

type Config struct {
	gorm.Model

	CTFName string
	StartAt time.Time
	EndAt   time.Time

	RegisterOpen bool
	CTFOpen      bool

	LockCount     int
	LockFrequency int
	LockDuration  int

	MinScore   uint
	MaxScore   uint
	CountToMin uint
}
