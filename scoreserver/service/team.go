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

type TeamApp interface {
	Login(teamname, password string) (*model.LoginToken, error)
	RegisterTeam(teamname, password, email string) (*model.Team, error)
	GetAdminTeam() (*model.Team, error)
	MakeTeamAdmin(t *model.Team) error
	GetTeamByID(teamID uint) (*model.Team, error)
	GetLoginTeam(token string) (*model.Team, error)
	PasswordResetRequest(email string) error
	PasswordReset(token, newpassword string) error
	PasswordUpdate(team *model.Team, newpassword string) error
}

func (app *app) RegisterTeam(teamname, password, email string) (*model.Team, error) {
	if err := app.validateTeamname(teamname); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if password == "" {
		return nil, xerrors.Errorf(": %w", NewErrorMessage("password is required"))
	}
	if err := app.validateEmail(email); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	t := model.Team{
		Teamname:     teamname,
		PasswordHash: hashPassword(password),
		Email:        email,
	}
	if err := app.repo.RegisterTeam(&t); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &t, nil
}

func (app *app) GetAdminTeam() (*model.Team, error) {
	t, err := app.repo.GetAdminTeam()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return t, nil
}

func (app *app) MakeTeamAdmin(t *model.Team) error {
	if err := app.repo.MakeTeamAdmin(t); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) GetTeamByID(teamID uint) (*model.Team, error) {
	team, err := app.repo.GetTeamByID(teamID)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (app *app) GetLoginTeam(token string) (*model.Team, error) {
	team, err := app.repo.GetTeamByLoginToken(token)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, NewErrorMessage("invalid token")
	} else if err != nil {
		return nil, err
	}
	return team, nil
}

func (app *app) Login(teamname, password string) (*model.LoginToken, error) {
	t, err := app.repo.GetTeamByName(teamname)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, NewErrorMessage("no such team")
	} else if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if !checkPassword(password, []byte(t.PasswordHash)) {
		return nil, NewErrorMessage("password mismatch")
	}

	token := model.LoginToken{
		TeamId:    t.ID,
		Token:     newToken(),
		ExpiresAt: tokenExpiredTime(),
	}
	if err := app.repo.SetTeamLoginToken(&token); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &token, nil
}

func (app *app) PasswordResetRequest(email string) error {
	t, err := app.repo.GetTeamByEmail(email)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return NewErrorMessage("invalid email")
		}
		return err
	}

	token := model.PasswordResetToken{
		TeamId:    t.ID,
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

	t, err := app.repo.GetTeamByPasswordResetToken(token)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return NewErrorMessage("token is invalid or expired")
		}
		return xerrors.Errorf(": %w", err)
	}

	if err := app.repo.UpdateTeamPassword(t, hashPassword(newpassword)); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// revoke *ALL* password reset token
	if err := app.repo.RevokeTeamPasswordResetToken(t.ID); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (app *app) PasswordUpdate(team *model.Team, newpassword string) error {
	if newpassword == "" {
		return NewErrorMessage("password is required")
	}
	if err := app.repo.UpdateTeamPassword(team, hashPassword(newpassword)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) validateTeamname(teamname string) error {
	if teamname == "" {
		return NewErrorMessage("teamname is required")
	}
	if len(teamname) >= 128 {
		return NewErrorMessage("teamname should be shorter than 128")
	}

	if _, err := app.repo.GetTeamByName(teamname); err == nil {
		return NewErrorMessage("teamname already used")
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

	if _, err := app.repo.GetTeamByEmail(email); err == nil {
		return NewErrorMessage("email already used")
	} else if err != nil && !xerrors.As(err, &repository.NotFoundError{}) {
		log.Printf("%v\n", err)
		log.Printf("%v\n", xerrors.As(err, &repository.NotFoundError{}))
		return xerrors.Errorf(": %w", err)
	}
	return nil
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
