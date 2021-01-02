package seeder

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
	"syreclabs.com/go/faker"
)

type Seeder interface {
	Team() (*model.Team, error)
	Challenge(flagFormat string) (*model.Challenge, error)
	Submission(t *model.Team, c *model.Challenge) (bool, error)
}

type seeder struct {
	app service.App
}

func New(app service.App) Seeder {
	return &seeder{
		app: app,
	}
}

func (s *seeder) Team() (*model.Team, error) {
	t, err := s.app.RegisterTeam(
		faker.Team().Name()+uuid.New().String(),
		faker.Internet().Password(10, 20),
		uuid.New().String()+faker.Internet().SafeEmail(),
		faker.Address().CountryCode(),
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return t, nil
}

func (s *seeder) Challenge(flagFormat string) (*model.Challenge, error) {
	flag := fmt.Sprintf(flagFormat, faker.Lorem().Word()+uuid.New().String())
	name := faker.Lorem().Sentence(3) + uuid.New().String()
	if err := s.app.AddChallenge(&service.Challenge{
		ID:          uuid.New().ID(),
		Name:        name,
		Flag:        flag,
		Description: faker.Lorem().String() + flag,
		Author:      faker.Name().LastName(),
		Score:       0,
		IsOpen:      true,
		IsRunning:   true,
		IsSurvey:    false,
	}); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	c, err := s.app.GetRawChallengeByName(name)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return c, nil
}

func (s *seeder) Submission(t *model.Team, c *model.Challenge) (bool, error) {
	var flag string
	if c != nil {
		flag = c.Flag
	} else {
		flag = faker.Hacker().IngVerb()
	}

	_, _, is_correct, err := s.app.SubmitFlag(t, faker.Internet().IpV4Address(), flag, true)
	if err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return is_correct, nil
}
