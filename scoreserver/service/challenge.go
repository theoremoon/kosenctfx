package service

import (
	"errors"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	_ "github.com/mattn/anko/packages"
	"github.com/mattn/anko/vm"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
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
	Compose     string       `json:"compose"`
	Deployment  string       `json:"deployment"`

	IsOpen   bool `json:"is_open"`
	IsSurvey bool `json:"is_survey"`
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
}

func (app *app) insertSubmission(s *model.Submission) error {
	if err := app.db.Create(s).Error; err != nil {
		return err
	}
	return nil
}

func (app *app) insertValidableSubmission(s *model.Submission) (bool, error) {
	valid := false

	// ややこしいのでtransactionをとって行う
	err := app.db.Transaction(func(tx *gorm.DB) error {
		// とりあえずsubmissionは保存しておく
		if err := app.db.Create(s).Error; err != nil {
			return err
		}

		// 既存の提出を読んでvalidityを決定する
		var count int64
		if err := app.db.Model(&model.ValidSubmission{}).Where("team_id = ? AND challenge_id = ?", s.TeamId, s.ChallengeId).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			valid = true
		}

		if valid {
			// validならvalid submissionを作成する
			vs := model.ValidSubmission{
				SubmissionId: s.ID,
				ChallengeId:  *s.ChallengeId,
				TeamId:       s.TeamId,
			}
			if err := app.db.Create(&vs).Error; err != nil {
				// ただしConstraint Errorが起きたらやはりValidではなかった
				if isDuplicatedError(err) {
					valid = false
					return nil
				} else {
					return xerrors.Errorf(": %w", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return valid, nil
}

func (app *app) markSubmissionValid(id uint32) error {
	if err := app.db.Model(&model.Submission{}).Where("id = ?", id).Update("is_valid", true).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) listTagsByChallengeIDs(ids []uint32) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := app.db.Order("challenge_id asc").Where("challenge_id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return tags, nil
}

func (app *app) listAttachmentsByChallengeIDs(ids []uint32) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	if err := app.db.Order("challenge_id asc").Where("challenge_id IN ?", ids).Find(&attachments).Error; err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return attachments, nil
}

func (app *app) listAllTags() ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := app.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (app *app) listAllAttachments() ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	if err := app.db.Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (app *app) rawChallengesToChallenges(cs []*model.Challenge) ([]*Challenge, error) {
	ids := make([]uint32, len(cs))
	for i, c := range cs {
		ids[i] = c.ID
	}

	tags, err := app.listTagsByChallengeIDs(ids)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	attachments, err := app.listAttachmentsByChallengeIDs(ids)
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
			Compose:     c.Compose,
			Deployment:  c.Deployment,
			Score:       0,            //TODO
			SolvedBy:    []SolvedBy{}, // TODO
			IsOpen:      c.IsOpen,
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
		Compose:     c.Compose,
		Deployment:  c.Deployment,
		Score:       0,            //TODO
		SolvedBy:    []SolvedBy{}, // TODO
		IsOpen:      c.IsOpen,
		IsSurvey:    c.IsSurvey,
	}

	tags, err := app.listTagsByChallengeIDs([]uint32{c.ID})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	chal.Tags = make([]string, len(tags))
	for i := range tags {
		chal.Tags[i] = tags[i].Tag
	}

	attachments, err := app.listAttachmentsByChallengeIDs([]uint32{c.ID})
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
	var challenges []*model.Challenge
	if err := app.db.Where("id IN ?", ids).Find(&challenges).Error; err != nil {
		return nil, err
	}
	chals, err := app.rawChallengesToChallenges(challenges)
	if err != nil {
		return nil, err
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

func (app *app) GetChallengeByFlag(flag string) (*model.Challenge, error) {
	var c model.Challenge
	if err := app.db.Where("flag = ?", flag).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (app *app) GetRawChallengeByID(challengeID uint32) (*model.Challenge, error) {
	var c model.Challenge
	if err := app.db.Where("id = ?", challengeID).First(&c).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewErrorMessage(challengeNotfoundMessage)
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &c, nil
}

func (app *app) GetRawChallengeByName(name string) (*model.Challenge, error) {
	var c model.Challenge
	if err := app.db.Where("name = ?", name).First(&c).Error; err != nil {
		if xerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewErrorMessage(challengeNotfoundMessage)
		}
		return nil, xerrors.Errorf(": %w", err)
	}
	return &c, nil
}

func (app *app) ListOpenedRawChallenges() ([]*model.Challenge, error) {
	var challenges []*model.Challenge
	if err := app.db.Where("is_open = ?", true).Find(&challenges).Error; err != nil {
		return nil, err
	}
	return challenges, nil
}

func (app *app) ListAllRawChallenges() ([]*model.Challenge, error) {
	var challenges []*model.Challenge
	if err := app.db.Find(&challenges).Error; err != nil {
		return nil, err
	}
	return challenges, nil
}

func (app *app) OpenChallenge(challengeID uint32) error {
	err := app.db.Model(&model.Challenge{}).
		Where("id = ?", challengeID).
		Update("is_open", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (app *app) CloseChallenge(challengeID uint32) error {
	err := app.db.Model(&model.Challenge{}).
		Where("id = ?", challengeID).
		Update("is_open", false).Error
	if err != nil {
		return err
	}
	return nil
}

func (app *app) addChallengeTag(t *model.Tag) error {
	if err := app.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (app *app) deleteTagByChallengeId(challengeId uint32) error {
	if err := app.db.Where("challenge_id = ?", challengeId).Delete(&model.Tag{}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) addChallengeAttachment(a *model.Attachment) error {
	if err := app.db.Create(a).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) deleteAttachmentByChallengeId(challengeId uint32) error {
	if err := app.db.Where("challenge_id = ?", challengeId).Delete(&model.Attachment{}).Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (app *app) AddChallenge(c *Challenge) error {
	if err := validateDeployment(c.Deployment); err != nil {
		return err
	}
	chal := model.Challenge{
		Name:        c.Name,
		Flag:        c.Flag,
		Category:    c.Category,
		Description: c.Description,
		Author:      c.Author,
		IsOpen:      false,
		IsSurvey:    c.IsSurvey,
		Compose:     c.Compose,
		Deployment:  c.Deployment,
	}
	if err := app.db.Create(&chal).Error; err != nil {
		if isDuplicatedError(err) {
			return NewErrorMessage(fmt.Sprintf(challengeDuplicatedMessage, c.Name))
		}
		return err
	}

	for _, t := range c.Tags {
		// do not care about error of this
		_ = app.addChallengeTag(&model.Tag{
			ChallengeId: chal.ID,
			Tag:         t,
		}).Error
	}

	for _, a := range c.Attachments {
		// do not care about error of this
		_ = app.addChallengeAttachment(&model.Attachment{
			ChallengeId: chal.ID,
			Name:        a.Name,
			URL:         a.URL,
		})
	}
	return nil
}

func (app *app) UpdateChallenge(challengeID uint32, c *Challenge) error {
	if err := validateDeployment(c.Deployment); err != nil {
		return err
	}

	chal := model.Challenge{
		Name:        c.Name,
		Flag:        c.Flag,
		Category:    c.Category,
		Description: c.Description,
		Author:      c.Author,
		IsSurvey:    c.IsSurvey,
		IsOpen:      c.IsOpen,
		Compose:     c.Compose,
		Deployment:  c.Deployment,
	}
	chal.ID = challengeID

	if err := app.db.Save(&chal).Error; err != nil {
		return err
	}

	// TODO: remove remote tags and attachments
	if err := app.deleteTagByChallengeId(challengeID); err != nil {
		return err
	}
	if err := app.deleteAttachmentByChallengeId(challengeID); err != nil {
		return err
	}

	for _, t := range c.Tags {
		// do not care about error of this
		_ = app.addChallengeTag(&model.Tag{
			ChallengeId: chal.ID,
			Tag:         t,
		})
	}

	for _, a := range c.Attachments {
		// do not care about error of this
		_ = app.addChallengeAttachment(&model.Attachment{
			ChallengeId: chal.ID,
			Name:        a.Name,
			URL:         a.URL,
		})
	}
	return nil
}

/// 返り値は 解いたchallenge（is_correctがfalseならnil)、 is_correct, is_valid, error
func (app *app) SubmitFlag(team *model.Team, ipaddress string, flag string, ctfRunning bool, submitted_at int64) (*model.Challenge, bool, bool, error) {
	chal, err := app.GetChallengeByFlag(flag)
	if err != nil && !xerrors.As(err, gorm.ErrRecordNotFound) {
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
		if err := app.insertSubmission(s); err != nil {
			return nil, false, false, xerrors.Errorf(": %w", err)
		}
		return nil, false, false, nil
	} else {
		// correct
		s.ChallengeId = &chal.ID
		s.IsCorrect = true

		if ctfRunning {
			// ctfRunningがtrueなときは初回の提出だけvalidになる。ここトランザクションかけておく
			valid, err := app.insertValidableSubmission(s)
			if err != nil {
				return nil, false, false, xerrors.Errorf(": %w", err)
			}

			if valid {
				// validなときsubmissionにvalidフラグ立てておく
				if err := app.markSubmissionValid(s.ID); err != nil {
					log.Errorf("%+v\n", err) // XXX
				}
			}
			return chal, true, valid, nil
		} else {
			// elseの場合は参考記録なのでvalidにしない
			if err := app.insertSubmission(s); err != nil {
				return nil, false, false, xerrors.Errorf(": %w", err)
			}
			return chal, true, false, nil
		}

	}
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

func validateDeployment(deploymentType string) error {
	if deploymentType == "" {
		return nil
	}
	if deploymentType == "one" {
		return nil
	}
	if deploymentType == "many" {
		return nil
	}
	return errors.New("deployment type should be one of `one`, `many`, ``")
}
