package seeder

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"math/rand"

	"github.com/google/uuid"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
	"syreclabs.com/go/faker"
)

//TODO: random wordにとを使う

var (
	TAGS = []string{"crypto", "web", "pwn", "misc", "cheating", "タグ"}
)

type Seeder interface {
	Team() (*model.Team, error)
	Challenge(flagFormat string) (*model.Challenge, error)
	Submission(t *model.Team, c *model.Challenge, submitted_at int64) (bool, error)
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
		randomname(),
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
	name := randomname()

	// 確率でsurvey問題を出す
	isSurvey := false
	if rand.Intn(10) == 0 {
		isSurvey = true
	}

	// tag
	tags := make([]string, 0)
	for _, t := range TAGS {
		if rand.Intn(10) == 0 {
			tags = append(tags, t)
		}
	}

	if err := s.app.AddChallenge(&service.Challenge{
		ID:          uuid.New().ID(),
		Name:        name,
		Flag:        flag,
		Description: faker.Lorem().String() + flag,
		Author:      faker.Name().LastName(),
		Score:       0,
		IsOpen:      true,
		IsRunning:   true,
		IsSurvey:    isSurvey,
		Tags:        tags,
	}); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	c, err := s.app.GetRawChallengeByName(name)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return c, nil
}

func (s *seeder) Submission(t *model.Team, c *model.Challenge, submitted_at int64) (bool, error) {
	var flag string
	if c != nil {
		flag = c.Flag
	} else {
		flag = faker.Hacker().IngVerb()
	}

	_, _, is_correct, err := s.app.SubmitFlag(t, faker.Internet().IpV4Address(), flag, true, submitted_at)
	if err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return is_correct, nil
}

func randomname() string {
	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	randomURL := "https://en.wikipedia.org/wiki/Special:Random"
	// 時々日本語になる
	if rand.Intn(10) == 0 {
		randomURL = "https://ja.wikipedia.org/wiki/%E7%89%B9%E5%88%A5:%E3%81%8A%E3%81%BE%E3%81%8B%E3%81%9B%E8%A1%A8%E7%A4%BA"
	}
	resp, err := client.Get(randomURL)
	if err != nil {
		panic(err)
	}
	words := strings.Split(resp.Header.Get("Location"), "/")
	word, err := url.QueryUnescape(words[len(words)-1])
	if err != nil {
		panic(err)
	}
	return word
}
