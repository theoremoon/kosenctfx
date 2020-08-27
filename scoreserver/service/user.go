package service

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	TeamID   uint   `json:"team_id"`
}

type UserApp interface {
	UserFilter(u *model.User) *User
	UsersFilter(us []*model.User) []*User

	LoginUser(username, password string) (*model.LoginToken, error)
	GetLoginUser(token string) (*model.User, error)
	GetUserByID(userID uint) (*model.User, error)
	GetTeamMembers(teamID uint) ([]*model.User, error)

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
	hpassword := sha256.Sum256([]byte(password))
	return bcrypt.CompareHashAndPassword(hashedPassword, hpassword[:]) == nil
}

func (app *app) UserFilter(u *model.User) *User {
	return &User{
		ID:       u.ID,
		Username: u.Username,
		TeamID:   u.TeamId,
	}
}

func (app *app) UsersFilter(us []*model.User) []*User {
	us2 := make([]*User, len(us))
	for i := 0; i < len(us); i++ {
		us2[i] = app.UserFilter(us[i])
	}
	return us2
}

func (app *app) LoginUser(username, password string) (*model.LoginToken, error) {
	u, err := app.repo.GetUserByUsername(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, NewErrorMessage("no such user")
	} else if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if !checkPassword(password, []byte(u.PasswordHash)) {
		return nil, NewErrorMessage("password mismatch")
	}

	token := model.LoginToken{
		UserId:    u.ID,
		Token:     newToken(),
		ExpiresAt: tokenExpiredTime(),
	}
	if err := app.repo.SetUserLoginToken(&token); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &token, nil
}

func (app *app) GetLoginUser(token string) (*model.User, error) {
	user, err := app.repo.GetUserByLoginToken(token)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, NewErrorMessage("invalid token")
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (app *app) GetUserByID(userID uint) (*model.User, error) {
	user, err := app.repo.GetUserByID(userID)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return user, nil
}

func (app *app) GetTeamMembers(teamID uint) ([]*model.User, error) {
	users, err := app.repo.GetTeamMembers(teamID)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return users, nil
}

func (app *app) RegisterUserWithTeam(username, password, email, teamname string) error {
	if err := app.validateUsername(username); err != nil {
		return err
	}
	if password == "" {
		return NewErrorMessage("password is required")
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
		return NewErrorMessage("password is required")
	}
	if err := app.validateEmail(email); err != nil {
		return err
	}

	t, err := app.repo.GetTeamByToken(teamToken)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return NewErrorMessage("invalid token")
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
	u, err := app.repo.GetUserByEmail(email)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return NewErrorMessage("invalid email")
		}
		return err
	}

	token := model.PasswordResetToken{
		UserId:    u.ID,
		Token:     newToken(),
		ExpiresAt: tokenExpiredTime(),
	}
	if err := app.repo.NewPasswordResetToken(&token); err != nil {
		return err
	}

	if err := app.mailer.Send(email, "password reset token", fmt.Sprintf("your password reset token is: %s", token.Token)); err != nil {
		return err
	}
	return nil
}

func (app *app) PasswordReset(token, newpassword string) error {
	if newpassword == "" {
		return NewErrorMessage("password is required")
	}

	u, err := app.repo.GetUserByPasswordResetToken(token)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return NewErrorMessage("token is invalid or expired")
		}
		return xerrors.Errorf(": %w", err)
	}

	if err := app.repo.UpdateUserPassword(u, hashPassword(newpassword)); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// revoke *ALL* password reset token
	if err := app.repo.RevokeUserPasswordResetToken(u.ID); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (app *app) PasswordUpdate(user *model.User, newpassword string) error {
	if newpassword == "" {
		return NewErrorMessage("password is required")
	}
	if err := app.repo.UpdateUserPassword(user, hashPassword(newpassword)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
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
		return NewErrorMessage("username is required")
	}
	if len(username) >= 128 {
		return NewErrorMessage("username should be shorter than 128")
	}

	if _, err := app.repo.GetUserByUsername(username); err == nil {
		return NewErrorMessage("username already used")
	} else if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	return nil
}

func (app *app) validateEmail(email string) error {
	if email == "" {
		return NewErrorMessage("email is required")
	}
	if len(email) >= 127 {
		return NewErrorMessage("email should be shorter than 128")
	}

	if _, err := app.repo.GetUserByEmail(email); err == nil {
		return NewErrorMessage("email already used")
	} else if err != nil && !xerrors.As(err, &repository.NotFoundError{}) {
		log.Printf("%v\n", err)
		log.Printf("%v\n", xerrors.As(err, &repository.NotFoundError{}))
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
