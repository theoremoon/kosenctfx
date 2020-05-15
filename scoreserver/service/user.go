package service

import (
	"crypto/sha256"

	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/crypto/bcrypt"
)

type UserApp interface {
	LoginUser(username, password string) (bool, error)
	LogoutUser(userId uint) error

	RegisterUserWithTeam(username, password, email, teamname string) error
	RegisterUserAndJoinToTeam(username, password, email, teamToken string) error
}

func hashPassword(password string) string {
	sha256password := sha256.Sum256([]byte(password))
	passwordHash, err := bcrypt.GenerateFromPassword(sha256password[:], bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(passwordHash)
}

func (app *app) LoginUser(username, password string) (bool, error) {
	u, err := app.repo.GetUserByUsername(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return false, ErrorMessage("no such user")
	} else if err != nil {
		return false, err
	}
	passwordHash := hashPassword(password)
	if u.PasswordHash != passwordHash {
		return false, ErrorMessage("password mismatch")
	}

	if err := app.repo.SetUserLoginToken(&model.LoginToken{
		UserId:    u.ID,
		Token:     newToken(),
		ExpiresAt: tokenExpiredTime(),
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (app *app) LogoutUser(userId uint) error {
	if err := app.repo.RevokeUserLoginToken(userId); err != nil {
		return err
	}
	return nil
}

func (app *app) RegisterUserWithTeam(username, password, email, teamname string) error {
	if err := app.validateUsername(username); err != nil {
		return err
	}
	if password == "" {
		return ErrorMessage("password is required")
	}
	if err := app.validateEmail(email); err != nil {
		return err
	}
	t, err := app.RegisterTeam(teamname)
	if err != nil {
		return err
	}
	u := model.User{
		Username:     username,
		PasswordHash: hashPassword(password),
		Email:        email,
		TeamId:       t.ID,
	}
	if err := app.repo.Register(&u); err != nil {
		return err
	}
	return nil
}

func (app *app) RegisterUserAndJoinToTeam(username, password, email, teamToken string) error {
	if err := app.validateUsername(username); err != nil {
		return err
	}
	if password == "" {
		return ErrorMessage("password is required")
	}
	if err := app.validateEmail(email); err != nil {
		return err
	}

	t, err := app.repo.GetTeamByToken(teamToken)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return ErrorMessage("invalid token")
	} else if err != nil {
		return err
	}

	u := model.User{
		Username:     username,
		PasswordHash: hashPassword(password),
		Email:        email,
		TeamId:       t.ID,
	}
	if err := app.repo.Register(&u); err != nil {
		return err
	}
	return nil

}

func (app *app) validateUsername(username string) error {
	if username == "" {
		return ErrorMessage("username is required")
	}
	if len(username) >= 128 {
		return ErrorMessage("username should be shorter than 128")
	}

	if _, err := app.repo.GetUserByUsername(username); err == nil {
		return ErrorMessage("username already used")
	} else if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	return nil
}

func (app *app) validateEmail(email string) error {
	if email == "" {
		return ErrorMessage("email is required")
	}
	if len(email) >= 127 {
		return ErrorMessage("email should be shorter than 128")
	}

	if _, err := app.repo.GetUserByEmail(email); err == nil {
		return ErrorMessage("email already used")
	} else if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	return nil
}
