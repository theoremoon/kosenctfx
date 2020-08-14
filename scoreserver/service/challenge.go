package service

import (
	"math"
	"path/filepath"

	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"golang.org/x/xerrors"
)

type Attachment struct {
	Name string
	URL  string
}

type Challenge struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Flag        string       `json:"flag"`
	Description string       `json:"description"`
	Author      string       `json:"author"`
	Score       uint         `json:"score"`
	Tags        []string     `json:"tags"`
	Attachments []Attachment `json:"attachments"`
	SolvedBy    []string     `json:"solved_by"`

	IsOpen    bool `json:"is_open"`
	IsRunning bool `json:"is_running"`
	IsSurvey  bool `json:"is_survey"`
}

type ChallengeApp interface {
	GetChallengeByID(challengeID uint) (*Challenge, error)
	GetChallengeByName(name string) (*Challenge, error)
	GetRawChallengeByName(name string) (*model.Challenge, error)
	ListAllChallenges() ([]*Challenge, error)
	ListOpenChallenges() ([]*Challenge, error)

	AddChallenge(c *Challenge) error
	OpenChallenge(challlengeID uint) error
	UpdateChallenge(challengeID uint, c *Challenge) (*model.Challenge, error)

	SubmitFlag(user *model.User, flag string) (*model.Challenge, bool, bool, error)
}

func (app *app) GetChallengeByID(challengeID uint) (*Challenge, error) {
	return nil, NewErrorMessage("not implemented")
}

func (app *app) GetChallengeByName(name string) (*Challenge, error) {
	return nil, NewErrorMessage("not implemented")
}

func (app *app) GetRawChallengeByName(name string) (*model.Challenge, error) {
	chal, err := app.repo.GetChallengeByName(name)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return nil, NewErrorMessage("No such challenge")
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return chal, nil
}

func (app *app) ListAllChallenges() ([]*Challenge, error) {
	// list all challenges, tags, and attachments
	chals, err := app.repo.ListAllChallenges()
	if err != nil {
		return nil, err
	}
	tags, err := app.repo.ListAllTags()
	if err != nil {
		return nil, err
	}
	attachments, err := app.repo.ListAllAttachments()
	if err != nil {
		return nil, err
	}

	// list valid submissions and its author team to calculate score
	submissions, err := app.repo.ListValidSubmissions()
	if err != nil {
		return nil, err
	}
	teams, err := app.repo.ListAllTeams()
	if err != nil {
		return nil, err
	}

	// make mapping as challenge id is the key
	tagMap := make(map[uint][]string)
	for _, c := range chals {
		tagMap[c.ID] = make([]string, 0)
	}
	for _, t := range tags {
		tagMap[t.ChallengeId] = append(tagMap[t.ChallengeId], t.Tag)
	}

	attachmentMap := make(map[uint][]Attachment)
	for _, c := range chals {
		attachmentMap[c.ID] = make([]Attachment, 0)
	}
	for _, a := range attachments {
		attachmentMap[a.ChallengeId] = append(attachmentMap[a.ChallengeId], Attachment{
			Name: filepath.Base(a.URL),
			URL:  a.URL,
		})
	}
	teamMap := make(map[uint]string)
	for _, t := range teams {
		teamMap[t.ID] = t.Teamname
	}

	// key: challlenge id, value: team name who solved this chal
	solvedByMap := make(map[uint][]string)
	for _, c := range chals {
		solvedByMap[c.ID] = make([]string, 0)
	}
	for _, s := range submissions {
		solvedByMap[*s.ChallengeId] = append(solvedByMap[*s.ChallengeId], teamMap[s.TeamId])
	}

	conf, err := app.repo.GetConfig()
	if err != nil {
		return nil, err
	}

	// make structure
	challenges := make([]*Challenge, len(chals))
	for i, c := range chals {
		challenges[i] = &Challenge{
			ID:          c.ID,
			Name:        c.Name,
			Flag:        c.Flag,
			Description: c.Description,
			Author:      c.Author,
			Score:       CalcChallengeScore(conf, uint(len(solvedByMap[c.ID]))),
			Tags:        tagMap[c.ID],
			Attachments: attachmentMap[c.ID],
			SolvedBy:    solvedByMap[c.ID],
			IsOpen:      c.IsOpen,
			IsSurvey:    c.IsSurvey,
		}
	}
	return challenges, nil
}

func (app *app) ListOpenChallenges() ([]*Challenge, error) {
	chals, err := app.ListAllChallenges()
	if err != nil {
		return nil, err
	}
	result := make([]*Challenge, 0, len(chals))
	for _, c := range chals {
		if c.IsOpen {
			result = append(result, c)
		}
	}
	return result, nil
}

func (app *app) AddChallenge(c *Challenge) error {
	chal := model.Challenge{
		Name:        c.Name,
		Flag:        c.Flag,
		Description: c.Description,
		Author:      c.Author,
		IsOpen:      false,
		IsSurvey:    c.IsSurvey,
	}
	err := app.repo.AddChallenge(&chal)
	if err != nil {
		if xerrors.As(err, &repository.DuplicatedError{}) {
			return NewErrorMessage("challenge is duplicated: " + chal.Name)
		}
		return xerrors.Errorf(": %w", err)
	}

	for _, t := range c.Tags {
		// do not care about error of this
		app.repo.AddChallengeTag(&model.Tag{
			ChallengeId: chal.ID,
			Tag:         t,
		})
	}

	for _, a := range c.Attachments {
		// do not care about error of this
		app.repo.AddChallengeAttachment(&model.Attachment{
			ChallengeId: chal.ID,
			Name:        a.Name,
			URL:         a.URL,
		})
	}
	return nil
}

func (app *app) OpenChallenge(challengeID uint) error {
	if err := app.repo.OpenChallengeByID(challengeID); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) UpdateChallenge(challengeID uint, c *Challenge) (*model.Challenge, error) {
	return nil, NewErrorMessage("not implemented")
}

func (app *app) SubmitFlag(user *model.User, flag string) (*model.Challenge, bool, bool, error) {
	correct := false
	valid := false
	chal, err := app.repo.GetChallengeByFlag(flag)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, false, false, err
	} else if err != nil {
		correct = true
	}
	if correct {
		_, err = app.repo.FindValidSubmission(user.TeamId, chal.ID)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return nil, false, false, err
		}
		if gorm.IsRecordNotFoundError(err) {
			valid = true
		}
	}
	if err := app.repo.InsertSubmission(&model.Submission{
		ChallengeId: &chal.ID,
		UserId:      user.ID,
		TeamId:      user.TeamId,
		IsCorrect:   correct,
		IsValid:     valid,
		Flag:        flag,
	}); err != nil {
		return nil, false, false, err
	}
	return chal, correct, valid, nil
}

func CalcChallengeScore(conf *model.Config, solveCount uint) uint {
	max := float64(conf.MaxScore)
	min := float64(conf.MinScore)

	a := max - min
	tx := (a - 1) / a
	t := 0.5 * math.Log((1+tx)/(1-tx)) // tanh^{-1}(tx)

	r := a / max
	y := max * math.Tanh(float64(solveCount)/(float64(conf.CountToMin)/t))
	return uint(r*(max-y) + min)
}
