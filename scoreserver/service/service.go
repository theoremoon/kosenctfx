package service

import (
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/theoremoon/kosenctfx/scoreserver/mailer"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

const (
	challengeNotfoundMessage         = "No such challenge"
	challengeDuplicatedMessage       = "Challenge %s exists"
	countrycodeInvalidMessage        = "Invalid country code (Not valid as ISO 3166-1 alpha-2)"
	countrycodeRequiredMessage       = "Country code is required"
	emailDuplicatedMessage           = "This email address is already used"
	emailNotfoundMessage             = "Invalid email address"
	emailRequiredMessage             = "Email is required"
	emailTooLongMessage              = "Maximum length of email is 128"
	passwordRequiredMessage          = "Password is required"
	passwordResetMailBody            = "Your password reset token is: %s"
	passwordResetMailTitle           = "Password Reset Token"
	passwordResetTokenInvalidMessage = "Password reset token is invalid"
	teamNotfoundMessage              = "No such team"
	teamnameDuplicatedMessage        = "This team name has already been taken"
	teamnameRequiredMessage          = "Team name is required"
	teamnameTooLongMessage           = "Maximum length of your team name is 128"
	tokenInvalidMessage              = "Invalid token"
	wrongPasswordMessage             = "Wrong password"
)

type App interface {
	TeamApp
	ChallengeApp
	CTFApp
	SubmissionApp
	ScoreFeed(chals []*model.Challenge, teams []*model.Team, submissions []*model.Submission) ([]*Challenge, []*ScoreFeedEntry, error)
	TaskSolves() (map[*model.Challenge]int64, error)
}

type app struct {
	repo   repository.Repository
	db     *gorm.DB
	mailer mailer.Mailer
}

func New(db *gorm.DB, mailer mailer.Mailer) App {
	return &app{
		mailer: mailer,
		db:     db,
	}
}

var LoginTokenLifeSpan = 7 * 24 * time.Hour // default is 1week

func tokenExpiredTime() time.Time {
	return time.Now().Add(LoginTokenLifeSpan)
}

func newToken() string {
	return uuid.New().String()
}

type TaskStat struct {
	Score    uint32 `json:"points"`
	SolvedAt int64  `json:"time"`
}

/// jsonの名前めちゃくちゃに見えるけどctftimeに沿ってるはず
type ScoreFeedEntry struct {
	Pos            int                  `json:"pos"`
	Teamname       string               `json:"team"`
	Country        string               `json:"country"`
	Score          int                  `json:"score"`
	TaskStats      map[string]*TaskStat `json:"taskStats"`
	TeamID         uint32               `json:"team_id"`
	LastSubmission int64                `json:"last_submission"`
}

func (app *app) ScoreFeed(chals []*model.Challenge, teams []*model.Team, submissions []*model.Submission) ([]*Challenge, []*ScoreFeedEntry, error) {
	conf, err := app.repo.GetConfig()
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	// List Tags and Attachments
	tags, err := app.repo.ListAllTags()
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}
	attachments, err := app.repo.ListAllAttachments()
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	// make mapping as challenge id is the key
	tagMap := make(map[uint32][]string)
	for _, c := range chals {
		tagMap[c.ID] = make([]string, 0)
	}
	for _, t := range tags {
		tagMap[t.ChallengeId] = append(tagMap[t.ChallengeId], t.Tag)
	}

	attachmentMap := make(map[uint32][]Attachment)
	for _, c := range chals {
		attachmentMap[c.ID] = make([]Attachment, 0)
	}
	for _, a := range attachments {
		attachmentMap[a.ChallengeId] = append(attachmentMap[a.ChallengeId], Attachment{
			Name: filepath.Base(a.URL),
			URL:  a.URL,
		})
	}
	teamMap := make(map[uint32]string)
	for _, t := range teams {
		teamMap[t.ID] = t.Teamname
	}

	// key: challlenge id, value: team name who solved this chal
	solvedByMap := make(map[uint32][]SolvedBy)
	for _, c := range chals {
		solvedByMap[c.ID] = make([]SolvedBy, 0)
	}
	for _, s := range submissions {
		solvedByMap[*s.ChallengeId] = append(solvedByMap[*s.ChallengeId], SolvedBy{
			TeamName: teamMap[s.TeamId],
			TeamID:   s.TeamId,
			SolvedAt: s.SubmittedAt,
		})
	}

	// make structure
	challenges := make([]*Challenge, len(chals))
	for i, c := range chals {
		score, err := CalcChallengeScore(int(len(solvedByMap[c.ID])), conf.ScoreExpr)
		if err != nil {
			return nil, nil, xerrors.Errorf(": %w", err)
		}
		challenges[i] = &Challenge{
			ID:          c.ID,
			Name:        c.Name,
			Flag:        c.Flag,
			Category:    c.Category,
			Description: c.Description,
			Author:      c.Author,
			Score:       uint32(score),
			Tags:        tagMap[c.ID],
			Attachments: attachmentMap[c.ID],
			SolvedBy:    solvedByMap[c.ID],
			IsOpen:      c.IsOpen,
			IsSurvey:    c.IsSurvey,
		}
	}

	// ----

	chalMap := make(map[uint32]*Challenge)
	for _, c := range challenges {
		chalMap[c.ID] = c
	}

	teamSubmissions := make(map[uint32][]*model.Submission)
	for _, t := range teams {
		teamSubmissions[t.ID] = make([]*model.Submission, 0)
	}
	for _, s := range submissions {
		teamSubmissions[s.TeamId] = append(teamSubmissions[s.TeamId], s)
	}

	// とりあえずエントリを作成する
	scoreFeed := make([]*ScoreFeedEntry, len(teams))
	for i := 0; i < len(teams); i++ {
		var score uint32 = 0
		taskStats := make(map[string]*TaskStat)
		var lastSubmission int64 = 0

		for _, s := range teamSubmissions[teams[i].ID] {
			c, exist := chalMap[*s.ChallengeId]
			if !exist {
				continue //?
			}
			score += c.Score
			solvedAt := s.SubmittedAt
			taskStats[c.Name] = &TaskStat{
				Score:    c.Score,
				SolvedAt: solvedAt,
			}
			if !c.IsSurvey && lastSubmission < solvedAt {
				lastSubmission = solvedAt
			}
		}

		scoreFeed[i] = &ScoreFeedEntry{
			Pos:            0,
			Teamname:       teams[i].Teamname,
			Country:        teams[i].CountryCode,
			TeamID:         teams[i].ID,
			Score:          int(score),
			TaskStats:      taskStats,
			LastSubmission: lastSubmission,
		}
	}

	// スコアと最終提出時刻でsort
	sort.Slice(scoreFeed, func(i, j int) bool {
		if scoreFeed[i].Score == scoreFeed[j].Score {
			return scoreFeed[i].LastSubmission < scoreFeed[j].LastSubmission
		}
		return scoreFeed[i].Score > scoreFeed[j].Score
	})

	// Posの値を埋める
	for i := 0; i < len(scoreFeed); i++ {
		scoreFeed[i].Pos = i + 1
		if i != 0 && scoreFeed[i].Score == scoreFeed[i-1].Score && scoreFeed[i].LastSubmission == scoreFeed[i-1].LastSubmission {
			scoreFeed[i].Pos = scoreFeed[i-1].Pos
		}
	}

	lastTeam := len(scoreFeed)
	// CTF開催からは0点のチームは表示しない
	if CalcCTFStatus(conf) != CTFNotStarted {
		for i := 0; i < len(scoreFeed); i++ {
			if scoreFeed[i].Score <= 0 {
				lastTeam = i
				break
			}
		}
	}

	return challenges, scoreFeed[:lastTeam], nil
}

/// どの問題が何回解かれたかを見る
func (app *app) TaskSolves() (map[*model.Challenge]int64, error) {
	chals, err := app.ListOpenedRawChallenges()
	if err != nil {
		return nil, err
	}
	chalMap := make(map[uint32]*model.Challenge)
	for _, c := range chals {
		chalMap[c.ID] = c
	}

	solves := make(map[*model.Challenge]int64)
	for _, c := range chals {
		solves[c] = 0
	}

	submissions, err := app.ListValidSubmissions()
	if err != nil {
		return nil, err
	}
	for _, s := range submissions {
		// valid submissionでnilということはなかろう
		solves[chalMap[*s.ChallengeId]]++
	}
	return solves, nil
}
