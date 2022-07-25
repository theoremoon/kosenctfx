package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
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
	ListAllTeams() ([]*model.Team, error)
	CountTeams() (int64, error)
	GetAdminTeam() (*model.Team, error)
	MakeTeamAdmin(t *model.Team) error
	GetTeamByID(teamID uint32) (*model.Team, error)
	GetTeamByName(teamName string) (*model.Team, error)
	GetLoginTeam(token string) (*model.Team, error)
	PasswordResetRequest(email string) error
	PasswordReset(token, newpassword string) error
	PasswordUpdate(team *model.Team, newpassword string) error
	UpdateTeamname(team *model.Team, newTeamname string) error
	UpdateEmail(team *model.Team, newEmail string) error
	UpdateCountry(team *model.Team, newCountryCode string) error
}

var (
	gountry_query *gountries.Query
)

func init() {
	gountry_query = gountries.New()
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
	if err := app.db.Create(t).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return &t, nil
}

func (app *app) listTeams(all bool) ([]*model.Team, error) {
	cond := make(map[string]interface{})
	if !all {
		cond["is_admin"] = false
	}

	var teams []*model.Team
	if err := app.db.Where("is_admin = ?", false).Find(&teams).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return teams, nil
}

func (app *app) ListTeams() ([]*model.Team, error) {
	return app.listTeams(false)
}

func (app *app) ListAllTeams() ([]*model.Team, error) {
	return app.listTeams(true)
}

func (app *app) CountTeams() (int64, error) {
	var count int64
	if err := app.db.Model(&model.Team{}).Where("is_admin = ?", false).Count(&count).Error; err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}

func (app *app) GetAdminTeam() (*model.Team, error) {

	var t model.Team
	if err := app.db.Where("is_admin = ?", true).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (app *app) MakeTeamAdmin(t *model.Team) error {
	if err := app.db.Model(t).Update("is_admin", true).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) GetTeamByID(teamID uint32) (*model.Team, error) {
	var t model.Team
	if err := app.db.Where("id = ?", teamID).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (app *app) GetTeamByName(teamName string) (*model.Team, error) {
	var t model.Team
	if err := app.db.Where("teamname = ?", teamName).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (app *app) getTeamByEmail(email string) (*model.Team, error) {
	var t model.Team
	if err := app.db.Where("email = ?", email).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (app *app) GetLoginTeam(token string) (*model.Team, error) {
	now := time.Now().Unix()
	var loginToken model.LoginToken
	if err := app.db.Where("token = ? AND expires_at > ?", token, now).First(&loginToken).Error; err != nil {
		return nil, NewErrorMessage(tokenInvalidMessage)
	}

	var t model.Team
	if err := app.db.Where("id = ?", loginToken.TeamId).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (app *app) Login(teamname, password, ipaddress string) (*model.LoginToken, error) {
	t, err := app.GetTeamByName(teamname)
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

	if err := app.db.Create(token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (app *app) PasswordResetRequest(email string) error {
	t, err := app.getTeamByEmail(email)
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

	if err := app.db.Create(token).Error; err != nil {
		return err
	}
	if err := app.mailer.Send(email, passwordResetMailTitle, fmt.Sprintf(passwordResetMailBody, token.Token)); err != nil {
		return err
	}
	return nil
}

func (app *app) getTeamByPasswordResetToken(token string) (*model.Team, error) {
	var t model.Team
	var resetToken model.PasswordResetToken
	now := time.Now().Unix()
	if err := app.db.Where("token = ? AND expires_at > ?", token, now).First(&resetToken).Error; err != nil {
		return nil, err
	}
	if err := app.db.Where("id = ?", resetToken.TeamId).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (app *app) PasswordReset(token, newpassword string) error {
	if newpassword == "" {
		return NewErrorMessage(passwordRequiredMessage)
	}

	t, err := app.getTeamByPasswordResetToken(token)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return NewErrorMessage(passwordResetTokenInvalidMessage)
		}
		return xerrors.Errorf(": %w", err)
	}

	if err := app.PasswordUpdate(t, hashPassword(newpassword)); err != nil {
		return err
	}
	// revoke *ALL* password reset token
	if err := app.db.Where("team_id = ?", t.ID).Delete(model.PasswordResetToken{}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (app *app) PasswordUpdate(team *model.Team, newpassword string) error {
	if newpassword == "" {
		return NewErrorMessage(passwordRequiredMessage)
	}
	if err := app.db.Model(team).Update("password_hash", hashPassword(newpassword)).Error; err != nil {
		return err
	}
	return nil
}

func (app *app) UpdateTeamname(team *model.Team, newTeamname string) error {
	if newTeamname == "" {
		return NewErrorMessage(teamnameRequiredMessage)
	}
	if err := app.db.Model(team).Update("teamname", newTeamname).Error; err != nil {
		if isDuplicatedError(err) {
			return errors.New("duplicated")
		}
		return err
	}
	return nil
}

func isDuplicatedError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if xerrors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1062 {
			return true
		}
	}
	return false
}

func (app *app) UpdateEmail(team *model.Team, newEmail string) error {
	if err := app.validateEmail(newEmail); err != nil {
		return err
	}
	if err := app.db.Model(team).Update("email", newEmail).Error; err != nil {
		if isDuplicatedError(err) {
			return errors.New("duplicated")
		}
		return err
	}
	return nil
}

func (app *app) UpdateCountry(team *model.Team, newCountryCode string) error {
	country, err := validateCountryCode(newCountryCode)
	if err != nil {
		return err
	}

	if err := app.db.Model(team).Update("country_code", country).Error; err != nil {
		return err
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

	if _, err := app.GetTeamByName(teamname); err == nil {
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

	if _, err := app.getTeamByEmail(email); err == nil {
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
	c, err := gountry_query.FindCountryByAlpha(countryCode)
	if err != nil {
		return "", NewErrorMessage(countrycodeInvalidMessage)
	}
	return c.Alpha2, nil
}
