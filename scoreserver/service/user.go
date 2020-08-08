package service

import (
	"crypto/sha256"

	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
}

type UserApp interface {
	LoginUser(username, password string) (*model.LoginToken, error)
	GetLoginUser(token string) (*model.User, error)
	GetUserByID(userID uint) (*User, error)

	RegisterUserWithTeam(username, password, email, teamname string) error
	RegisterUserAndJoinToTeam(username, password, email, teamToken string) error

	PasswordResetRequest(email string) error
	PasswordReset(token, newpassword string) error
	PasswordUpdate(user *model.User, newpassword string) error

	GetAdminUser() (*model.User, error)
	CreateAdminUser(email, password string) (*model.User, error)
}

func hashPassword(password string) string {
	sha256password := sha256.Sum256([]byte(password))
	passwordHash, err := bcrypt.GenerateFromPassword(sha256password[:], bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(passwordHash)
}
func checkPassword(password string, hashedPassword []byte) bool {
	hpassword := hashPassword(password)
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(hpassword)) == nil
}

func (app *app) LoginUser(username, password string) (*model.LoginToken, error) {
	u, err := app.repo.GetUserByUsername(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, ErrorMessage("no such user")
	} else if err != nil {
		return nil, err
	}
	if checkPassword(password, []byte(u.PasswordHash)) {
		return nil, ErrorMessage("password mismatch")
	}

	token := model.LoginToken{
		UserId:    u.ID,
		Token:     newToken(),
		ExpiresAt: tokenExpiredTime(),
	}
	if err := app.repo.SetUserLoginToken(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (app *app) GetLoginUser(token string) (*model.User, error) {
	user, err := app.repo.GetUserByLoginToken(token)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, ErrorMessage("invalid token")
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (app *app) GetUserByID(userID uint) (*User, error) {
	return nil, ErrorMessage("not implemented")
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

func (app *app) PasswordResetRequest(email string) error {
	return ErrorMessage("not implemented")
}

func (app *app) PasswordReset(token, newpassword string) error {
	return ErrorMessage("not implemented")
}

func (app *app) PasswordUpdate(user *model.User, newpassword string) error {
	return ErrorMessage("not implemented")
}

func (app *app) GetAdminUser() (*model.User, error) {
	user, err := app.repo.GetAdminUser()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (app *app) CreateAdminUser(email, password string) (*model.User, error) {
	t, err := app.RegisterAdminTeam("admin")
	if err != nil {
		return nil, err
	}
	u := model.User{
		Username:     "admin",
		PasswordHash: hashPassword(password),
		Email:        email,
		TeamId:       t.ID,
		IsAdmin:      true,
	}
	if err := app.repo.Register(&u); err != nil {
		return nil, err
	}
	return &u, nil
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
