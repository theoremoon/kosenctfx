package service

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	_ "github.com/mattn/anko/packages"
	"github.com/mattn/anko/vm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"golang.org/x/xerrors"
)

type Attachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type SolvedBy struct {
	SolvedAt int64  `json:"solved_at"`
	TeamID   uint32 `json:"team_id"`
	TeamName string `json:"team_name"`
}

type Challenge struct {
	ID          uint32       `json:"id"`
	Name        string       `json:"name"`
	Flag        string       `json:"flag"`
	Category    string       `json:"category"`
	Description string       `json:"description"`
	Author      string       `json:"author"`
	Score       uint32       `json:"score"`
	Tags        []string     `json:"tags"`
	Attachments []Attachment `json:"attachments"`
	SolvedBy    []SolvedBy   `json:"solved_by"`

	IsOpen    bool `json:"is_open"`
	IsRunning bool `json:"is_running"`
	IsSurvey  bool `json:"is_survey"`
}

type ChallengeApp interface {
	GetChallengeByID(challengeID uint32) (*Challenge, error)
	ListChallengeByIDs(ids []uint32) ([]*Challenge, error)
	GetChallengeByName(name string) (*Challenge, error)
	GetRawChallengeByID(challengeID uint32) (*model.Challenge, error)
	GetRawChallengeByName(name string) (*model.Challenge, error)
	ListOpenedRawChallenges() ([]*model.Challenge, error)
	ListAllRawChallenges() ([]*model.Challenge, error)

	AddChallenge(c *Challenge) error
	OpenChallenge(challengeID uint32) error
	CloseChallenge(challengeID uint32) error
	UpdateChallenge(challengeID uint32, c *Challenge) error

	SubmitFlag(team *model.Team, ipaddress string, flag string, ctfRunning bool, submitted_at int64) (*model.Challenge, bool, bool, error)

	GetWrongCount(teamID uint32, duration time.Duration) (int64, error)
	LockSubmission(teamID uint32, duration time.Duration) error
	CheckSubmittable(teamID uint32) (bool, error)
}

func (app *app) rawChallengesToChallenges(cs []*model.Challenge) ([]*Challenge, error) {
	ids := make([]uint32, len(cs))
	for i, c := range cs {
		ids[i] = c.ID
	}

	tags, err := app.repo.ListTagsByChallengeIDs(ids)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	attachments, err := app.repo.ListAttachmentsByChallengeIDs(ids)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	tagMap := make(map[uint32][]string)
	for _, id := range ids {
		tagMap[id] = make([]string, 0)
	}
	for _, t := range tags {
		tagMap[t.ChallengeId] = append(tagMap[t.ChallengeId], t.Tag)
	}

	attachmentMap := make(map[uint32][]Attachment)
	for _, id := range ids {
		attachmentMap[id] = make([]Attachment, 0)
	}
	for _, a := range attachments {
		attachmentMap[a.ChallengeId] = append(attachmentMap[a.ChallengeId], Attachment{
			Name: a.Name,
			URL:  a.URL,
		})
	}

	chals := make([]*Challenge, len(cs))
	for i, c := range cs {
		chals[i] = &Challenge{
			ID:          c.ID,
			Name:        c.Name,
			Flag:        c.Flag,
			Description: c.Description,
			Author:      c.Author,
			Score:       0,            //TODO
			SolvedBy:    []SolvedBy{}, // TODO
			IsOpen:      c.IsOpen,
			IsRunning:   false, // TODO
			IsSurvey:    c.IsSurvey,
			Tags:        tagMap[c.ID],
			Attachments: attachmentMap[c.ID],
		}
	}
	return chals, nil
}

func (app *app) rawChallengeToChallenge(c *model.Challenge) (*Challenge, error) {
	chal := Challenge{
		ID:          c.ID,
		Name:        c.Name,
		Flag:        c.Flag,
		Description: c.Description,
		Author:      c.Author,
		Score:       0,            //TODO
		SolvedBy:    []SolvedBy{}, // TODO
		IsOpen:      c.IsOpen,
		IsRunning:   false, // TODO
		IsSurvey:    c.IsSurvey,
	}

	tags, err := app.repo.FindTagsByChallengeID(c.ID)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	chal.Tags = make([]string, len(tags))
	for i := range tags {
		chal.Tags[i] = tags[i].Tag
	}

	attachments, err := app.repo.FindAttachmentsByChallengeID(c.ID)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	chal.Attachments = make([]Attachment, len(attachments))
	for i := range attachments {
		chal.Attachments[i] = Attachment{
			Name: attachments[i].Name,
			URL:  attachments[i].URL,
		}
	}

	// TODO
	return &chal, nil
}

func (app *app) GetChallengeByID(challengeID uint32) (*Challenge, error) {
	c, err := app.GetRawChallengeByID(challengeID)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return app.rawChallengeToChallenge(c)
}

func (app *app) ListChallengeByIDs(ids []uint32) ([]*Challenge, error) {
	cs, err := app.repo.ListChallengeByIDs(ids)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	chals, err := app.rawChallengesToChallenges(cs)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return chals, nil
}

func (app *app) GetChallengeByName(name string) (*Challenge, error) {
	c, err := app.GetRawChallengeByName(name)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return app.rawChallengeToChallenge(c)
}

func (app *app) GetRawChallengeByID(challengeID uint32) (*model.Challenge, error) {
	chal, err := app.repo.GetChallengeByID(challengeID)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return nil, NewErrorMessage(challengeNotfoundMessage)
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return chal, nil
}

func (app *app) GetRawChallengeByName(name string) (*model.Challenge, error) {
	chal, err := app.repo.GetChallengeByName(name)
	if err != nil {
		if xerrors.As(err, &repository.NotFoundError{}) {
			return nil, NewErrorMessage(challengeNotfoundMessage)
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return chal, nil
}

func (app *app) ListOpenedRawChallenges() ([]*model.Challenge, error) {
	chals, err := app.repo.ListOpenedChallenges()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return chals, nil
}

func (app *app) ListAllRawChallenges() ([]*model.Challenge, error) {
	chals, err := app.repo.ListAllChallenges()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return chals, nil
}

func (app *app) AddChallenge(c *Challenge) error {
	chal := model.Challenge{
		Name:        c.Name,
		Flag:        c.Flag,
		Category:    c.Category,
		Description: c.Description,
		Author:      c.Author,
		IsOpen:      false,
		IsSurvey:    c.IsSurvey,
	}
	err := app.repo.AddChallenge(&chal)
	if err != nil {
		if xerrors.As(err, &repository.DuplicatedError{}) {
			return NewErrorMessage(fmt.Sprintf(challengeDuplicatedMessage, c.Name))
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

func (app *app) OpenChallenge(challengeID uint32) error {
	if err := app.repo.OpenChallengeByID(challengeID); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
func (app *app) CloseChallenge(challengeID uint32) error {
	if err := app.repo.CloseChallengeByID(challengeID); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil

}
func (app *app) UpdateChallenge(challengeID uint32, c *Challenge) error {
	chal := model.Challenge{
		Name:        c.Name,
		Flag:        c.Flag,
		Category:    c.Category,
		Description: c.Description,
		Author:      c.Author,
		IsSurvey:    c.IsSurvey,
		IsOpen:      c.IsOpen,
	}
	chal.ID = challengeID

	err := app.repo.UpdateChallenge(&chal)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if err := app.repo.DeleteTagByChallengeId(challengeID); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := app.repo.DeleteAttachmentByChallengeId(challengeID); err != nil {
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

/// 返り値は 解いたchallenge（is_correctがfalseならnil)、 is_correct, is_valid, error
func (app *app) SubmitFlag(team *model.Team, ipaddress string, flag string, ctfRunning bool, submitted_at int64) (*model.Challenge, bool, bool, error) {
	chal, err := app.repo.GetChallengeByFlag(flag)
	if err != nil && !xerrors.As(err, &repository.NotFoundError{}) {
		return nil, false, false, xerrors.Errorf(": %w", err)
	}

	s := &model.Submission{
		TeamId:      team.ID,
		IsCorrect:   false, //とりあえずfalseを入れておいてあとからtrueで上書きする
		IsValid:     false, // とりあえずfalseを入れておいてあとからtrueで上書きする
		Flag:        flag,
		IPAddress:   ipaddress,
		SubmittedAt: submitted_at,
	}
	if chal == nil || !chal.IsOpen {
		// wrong
		if err := app.repo.InsertSubmission(s); err != nil {
			return nil, false, false, xerrors.Errorf(": %w", err)
		}
		return nil, false, false, nil
	} else {
		// correct
		s.ChallengeId = &chal.ID
		s.IsCorrect = true

		if ctfRunning {
			// ctfRunningがtrueなときは初回の提出だけvalidになる。ここトランザクションかけておく
			valid, err := app.repo.InsertValidableSubmission(s)
			if err != nil {
				return nil, false, false, xerrors.Errorf(": %w", err)
			}

			if valid {
				// validなときsubmissionにvalidフラグ立てておく
				if err := app.repo.MarkSubmissionValid(s.ID); err != nil {
					log.Errorf("%+v\n", err) // XXX
				}
			}
			return chal, true, valid, nil
		} else {
			// elseの場合は参考記録なのでvalidにしない
			if err := app.repo.InsertSubmission(s); err != nil {
				return nil, false, false, xerrors.Errorf(": %w", err)
			}
			return chal, true, false, nil
		}

	}
}

func (app *app) GetWrongCount(teamID uint32, duration time.Duration) (int64, error) {
	cnt, err := app.repo.GetWrongCount(teamID, duration)
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}
	return cnt, nil
}

func (app *app) LockSubmission(teamID uint32, duration time.Duration) error {
	err := app.repo.LockSubmission(teamID, duration)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) CheckSubmittable(teamID uint32) (bool, error) {
	b, err := app.repo.CheckSubmittable(teamID)
	if err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return b, nil
}

func CalcChallengeScore(solveCount int, scoreExpr string) (int, error) {
	script := fmt.Sprintf(`
	var math = import("math")
	%s
	toInt(calc(count))
	`, scoreExpr)

	e := env.NewEnv()
	core.ImportToX(e)
	e.Define("count", solveCount)
	r, err := vm.Execute(e, nil, script)
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	switch v := r.(type) {
	case int:
		return v, nil
	case uint:
		return int(v), nil
	case int32:
		return int(v), nil
	case uint32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint64:
		return int(v), nil
	default:
		return 0, xerrors.Errorf("score calculation returns invalid type: %T", r)
	}
}
