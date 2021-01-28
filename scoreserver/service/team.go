package service

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/pariz/gountries"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type TeamApp interface {
	Login(teamname, password, ipaddress string) (*model.LoginToken, error)
	RegisterTeam(teamname, password, email, countryCode string) (*model.Team, error)
	ListTeams() ([]*model.Team, error)
	GetAdminTeam() (*model.Team, error)
	MakeTeamAdmin(t *model.Team) error
	GetTeamByID(teamID uint32) (*model.Team, error)
	GetLoginTeam(token string) (*model.Team, error)
	PasswordResetRequest(email string) error
	PasswordReset(token, newpassword string) error
	PasswordUpdate(team *model.Team, newpassword string) error
	UpdateTeamname(team *model.Team, newTeamname string) error
	UpdateCountry(team *model.Team, newCountryCode string) error
}

func (app *app) RegisterTeam(teamname, password, email, countryCode string) (*model.Team, error) {
	if err := app.validateTeamname(teamname); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if password == "" {
		return nil, xerrors.Errorf(": %w", NewErrorMessage(passwordRequiredMessage))
	}
	if err := app.validateEmail(email); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	country, err := validateCountryCode(countryCode)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	t := model.Team{
		Teamname:     teamname,
		PasswordHash: hashPassword(password),
		Email:        email,
		CountryCode:  country,
	}
	if err := app.repo.RegisterTeam(&t); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &t, nil
}

func (app *app) ListTeams() ([]*model.Team, error) {
	teams, err := app.repo.ListAllTeams()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return teams, nil
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

func (app *app) GetTeamByID(teamID uint32) (*model.Team, error) {
	team, err := app.repo.GetTeamByID(teamID)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (app *app) GetLoginTeam(token string) (*model.Team, error) {
	team, err := app.repo.GetTeamByLoginToken(token)
	if err != nil && xerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, NewErrorMessage(tokenInvalidMessage)
	} else if err != nil {
		return nil, err
	}
	return team, nil
}

func (app *app) Login(teamname, password, ipaddress string) (*model.LoginToken, error) {
	t, err := app.repo.GetTeamByName(teamname)
	if err != nil && xerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, NewErrorMessage(teamNotfoundMessage)
	} else if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if !checkPassword(password, []byte(t.PasswordHash)) {
		return nil, NewErrorMessage(wrongPasswordMessage)
	}

	token := model.LoginToken{
		TeamId:    t.ID,
		Token:     newToken(),
		ExpiresAt: tokenExpiredTime().Unix(),
		IPAddress: ipaddress,
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
			return NewErrorMessage(emailNotfoundMessage)
		}
		return err
	}

	token := model.PasswordResetToken{
		TeamId:    t.ID,
		Token:     newToken(),
		ExpiresAt: tokenExpiredTime().Unix(),
	}
	if err := app.repo.NewPasswordResetToken(&token); err != nil {
		return err
	}

	if err := app.mailer.Send(email, passwordResetMailTitle, fmt.Sprintf(passwordResetMailBody, token.Token)); err != nil {
		return err
	}
	return nil
}

func (app *app) PasswordReset(token, newpassword string) error {
	if newpassword == "" {
		return NewErrorMessage(passwordRequiredMessage)
	}

	t, err := app.repo.GetTeamByPasswordResetToken(token)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return NewErrorMessage(passwordResetTokenInvalidMessage)
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
		return NewErrorMessage(passwordRequiredMessage)
	}
	if err := app.repo.UpdateTeamPassword(team, hashPassword(newpassword)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) UpdateTeamname(team *model.Team, newTeamname string) error {
	if newTeamname == "" {
		return NewErrorMessage(teamnameRequiredMessage)
	}
	if err := app.repo.UpdateTeamname(team, newTeamname); err != nil {
		if xerrors.As(err, &repository.DuplicatedError{}) {
			return xerrors.Errorf(": %w", NewErrorMessage("that teamname is already used"))
		}
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) UpdateCountry(team *model.Team, newCountryCode string) error {
	if newCountryCode == "" {
		return NewErrorMessage(countrycodeRequiredMessage)
	}

	country, err := validateCountryCode(newCountryCode)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if err := app.repo.UpdateCountry(team, country); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) validateTeamname(teamname string) error {
	if teamname == "" {
		return NewErrorMessage(teamnameRequiredMessage)
	}
	if len(teamname) >= 128 {
		return NewErrorMessage(teamnameTooLongMessage)
	}

	if _, err := app.repo.GetTeamByName(teamname); err == nil {
		return NewErrorMessage(teamnameDuplicatedMessage)
	} else if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (app *app) validateEmail(email string) error {
	if email == "" {
		return NewErrorMessage(emailRequiredMessage)
	}
	if len(email) >= 127 {
		return NewErrorMessage(emailTooLongMessage)
	}

	if _, err := app.repo.GetTeamByEmail(email); err == nil {
		return NewErrorMessage(emailDuplicatedMessage)
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

func validateCountryCode(countryCode string) (string, error) {
	if countryCode == "" {
		return "", nil
	}
	q := gountries.New()
	c, err := q.FindCountryByAlpha(countryCode)
	if err != nil {
		return "", NewErrorMessage(countrycodeInvalidMessage)
	}
	return c.Alpha2, nil
}
