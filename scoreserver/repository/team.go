package repository

import (
	"time"

	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type TeamRepository interface {
	RegisterTeam(t *model.Team) error
	MakeTeamAdmin(t *model.Team) error
	GetAdminTeam() (*model.Team, error)
	ListAllTeams() ([]*model.Team, error)
	CountTeams() (int64, error)
	AdminListAllTeams() ([]*model.Team, error)
	GetTeamByID(teamId uint32) (*model.Team, error)
	GetTeamByLoginToken(token string) (*model.Team, error)
	GetTeamByName(teamname string) (*model.Team, error)
	GetTeamByEmail(email string) (*model.Team, error)
	GetTeamByPasswordResetToken(token string) (*model.Team, error)
	SetTeamLoginToken(token *model.LoginToken) error
	NewPasswordResetToken(token *model.PasswordResetToken) error
	RevokeTeamPasswordResetToken(teamID uint32) error
	UpdateTeamPassword(team *model.Team, passwordHash string) error
	UpdateTeamname(team *model.Team, teamname string) error
	UpdateEmail(team *model.Team, email string) error
	UpdateCountry(team *model.Team, country string) error
	UpdatePrizeEligible(team *model.Team, eligible bool) error
}

func (r *repository) RegisterTeam(t *model.Team) error {
	if err := r.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) MakeTeamAdmin(team *model.Team) error {
	if err := r.db.Model(team).Update("is_admin", true).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) GetAdminTeam() (*model.Team, error) {
	var t model.Team
	if err := r.db.Where("is_admin = ?", true).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) ListAllTeams() ([]*model.Team, error) {
	var teams []*model.Team
	if err := r.db.Where("is_admin = ?", false).Find(&teams).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return teams, nil
}

func (r *repository) CountTeams() (int64, error) {
	var count int64
	if err := r.db.Model(&model.Team{}).Where("is_admin = ?", false).Count(&count).Error; err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return count, nil
}

func (r *repository) AdminListAllTeams() ([]*model.Team, error) {
	var teams []*model.Team
	if err := r.db.Find(&teams).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return teams, nil
}

func (r *repository) GetTeamByID(teamID uint32) (*model.Team, error) {
	var t model.Team
	if err := r.db.Where("id = ?", teamID).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) GetTeamByLoginToken(token string) (*model.Team, error) {
	var t model.Team
	var loginToken model.LoginToken
	now := time.Now().Unix()
	if err := r.db.Where("token = ? AND expires_at > ?", token, now).First(&loginToken).Error; err != nil {
		return nil, err
	}
	if err := r.db.Where("id = ?", loginToken.TeamId).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) GetTeamByName(teamname string) (*model.Team, error) {
	var t model.Team
	if err := r.db.Where("teamname = ?", teamname).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) GetTeamByEmail(email string) (*model.Team, error) {
	var t model.Team
	if err := r.db.Where("email = ?", email).First(&t).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf(": %w", NotFound("team"))
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &t, nil
}

func (r *repository) GetTeamByPasswordResetToken(token string) (*model.Team, error) {
	var t model.Team
	var resetToken model.PasswordResetToken
	now := time.Now().Unix()
	if err := r.db.Where("token = ? AND expires_at > ?", token, now).First(&resetToken).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFound("token")
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	if err := r.db.Where("id = ?", resetToken.TeamId).First(&t).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFound("team")
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &t, nil
}

func (r *repository) SetTeamLoginToken(token *model.LoginToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) NewPasswordResetToken(token *model.PasswordResetToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) RevokeTeamPasswordResetToken(teamID uint32) error {
	if err := r.db.Where("team_id = ?", teamID).Delete(model.PasswordResetToken{}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) UpdateTeamPassword(team *model.Team, passwordHash string) error {
	if err := r.db.Model(team).Update("password_hash", passwordHash).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) UpdateTeamname(team *model.Team, teamname string) error {
	if err := r.db.Model(team).Update("teamname", teamname).Error; err != nil {
		if isDuplicatedError(err) {
			return xerrors.Errorf(": %w", Duplicated("teamname"))
		}
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) UpdateEmail(team *model.Team, email string) error {
	if err := r.db.Model(team).Update("email", email).Error; err != nil {
		if isDuplicatedError(err) {
			return xerrors.Errorf(": %w", Duplicated("email"))
		}
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) UpdateCountry(team *model.Team, country string) error {
	if err := r.db.Model(team).Update("country_code", country).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (r *repository) UpdatePrizeEligible(team *model.Team, eligible bool) error {
	if err := r.db.Model(team).Update("is_prize_eligible", eligible).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
