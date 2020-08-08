package repository

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type UserRepository interface {
	Register(u *model.User) error
	GetAdminUser() (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(username string) (*model.User, error)
	GetUserById(userId uint) (*model.User, error)
	GetUserByLoginToken(token string) (*model.User, error)
	SetUserLoginToken(token *model.LoginToken) error
}

func (r *repository) Register(u *model.User) error {
	err := r.db.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAdminUser() (*model.User, error) {
	var u model.User
	if err := r.db.Where("is_admin = ?", true).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetUserByUsername(username string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetUserByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetUserById(userId uint) (*model.User, error) {
	var u model.User
	if err := r.db.Where("id = ?", userId).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetUserByLoginToken(token string) (*model.User, error) {
	var u model.User
	var t model.LoginToken
	now := time.Now().Unix()
	if err := r.db.Where("token = ? AND expires_at > ?", token, now).First(&t).Error; err != nil {
		return nil, err
	}
	if err := r.db.Where("id = ?", t.UserId).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) SetUserLoginToken(token *model.LoginToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return err
	}
	return nil
}
