package service

import (
	"math"
	"path/filepath"

	"github.com/jinzhu/gorm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

type Attachment struct {
	Name string
	URL  string
}

type Challenge struct {
	Name        string
	Flag        string
	Description string
	Author      string
	Score       uint
	Tags        []string
	Attachments []Attachment
	SolvedBy    []string

	IsOpen   bool
	IsSurvey bool
}

type ChallengeApp interface {
	ListAllChallenges() ([]*Challenge, error)
	ListOpenChallenges() ([]*Challenge, error)
	SubmitFlag(user *model.User, flag string) (*model.Challenge, bool, bool, error)
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
	for _, t := range tags {
		tagMap[t.ChallengeId] = make([]string, 0)
	}
	for _, t := range tags {
		tagMap[t.ChallengeId] = append(tagMap[t.ChallengeId], t.Tag)
	}
	attachmentMap := make(map[uint][]Attachment)
	for _, a := range attachments {
		attachmentMap[a.ChallengeId] = make([]Attachment, 0)
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
	solvedByMap := make(map[uint][]string)
	for _, s := range submissions {
		solvedByMap[*s.ChallengeId] = make([]string, 0)
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
func (app *app) SubmitFlag(user *model.User, flag string) (*model.Challenge, bool, bool, error) {
	correct := false
	valid := false
	chal, err := app.repo.FindChallengeByFlag(flag)
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
