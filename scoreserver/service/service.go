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
)

type App interface {
	UserApp
	TeamApp
	ChallengeApp
	CTFApp
	NotificationApp
	ScoreFeed() ([]*Challenge, *Scoreboard, *Scoreboard, error)
}

type app struct {
	repo   repository.Repository
	mailer mailer.Mailer
}

func New(repo repository.Repository, mailer mailer.Mailer) App {
	return &app{
		mailer: mailer,
		repo:   repo,
	}
}

var LoginTokenLifeSpan = 7 * 24 * time.Hour // default is 1week

func tokenExpiredTime() time.Time {
	return time.Now().Add(LoginTokenLifeSpan)
}

func newToken() string {
	return uuid.New().String()
}

func (app *app) ValidateRunning(t time.Time) error {
	conf, err := app.repo.GetConfig()
	if err != nil {
		return err
	}
	if !t.After(conf.StartAt) {
		return NewErrorMessage("CTF has not started yet")
	}
	if !t.Before(conf.EndAt) {
		return NewErrorMessage("CTF has alredy finished")
	}
	return nil
}

type TaskStat struct {
	Score    uint  `json:"points"`
	SolvedAt int64 `json:"time"`
}
type ScoreFeedEntry struct {
	Pos            int                  `json:"pos"`
	Teamname       string               `json:"team"`
	Score          int                  `json:"points"`
	TaskStats      map[string]*TaskStat `json:"taskStats"`
	TeamID         uint                 `json:"team_id"`
	LastSubmission int64                `json:"last_submission"`
}
type Scoreboard struct {
	Tasks     []string          `json:"tasks"`
	Standings []*ScoreFeedEntry `json:"standings"`
}

/// 問題一覧、とチームのランキングを同時に計算する
func (app *app) ScoreFeed() ([]*Challenge, *Scoreboard, *Scoreboard, error) {
	conf, err := app.repo.GetConfig()
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
	}
	status := CalcCTFStatus(conf)

	// list all challenges, tags, and attachments
	allchals, err := app.repo.ListAllChallenges()
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
	}
	chals := make([]*model.Challenge, 0, len(allchals))
	if status == CTFRunning || status == CTFEnded {
		for _, c := range allchals {
			if c.IsOpen {
				chals = append(chals, c)
			}
		}
	}

	tags, err := app.repo.ListAllTags()
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
	}
	attachments, err := app.repo.ListAllAttachments()
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
	}

	// list valid submissions and its author team to calculate score
	submissions, err := app.repo.ListValidSubmissions()
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
	}
	teams, err := app.repo.ListAllTeams()
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
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
	solvedByMap := make(map[uint][]SolvedBy)
	for _, c := range chals {
		solvedByMap[c.ID] = make([]SolvedBy, 0)
	}
	for _, s := range submissions {
		solvedByMap[s.ChallengeId] = append(solvedByMap[s.ChallengeId], SolvedBy{
			TeamName: teamMap[s.TeamId],
			TeamID:   s.TeamId,
			SolvedAt: s.CreatedAt.Unix(),
		})
	}

	// make structure
	challenges := make([]*Challenge, len(chals))
	for i, c := range chals {
		score, err := CalcChallengeScore(int(len(solvedByMap[c.ID])), conf.ScoreExpr)
		if err != nil {
			return nil, nil, nil, xerrors.Errorf(": %w", err)
		}
		challenges[i] = &Challenge{
			ID:          c.ID,
			Name:        c.Name,
			Flag:        c.Flag,
			Description: c.Description,
			Author:      c.Author,
			Score:       uint(score),
			Tags:        tagMap[c.ID],
			Attachments: attachmentMap[c.ID],
			SolvedBy:    solvedByMap[c.ID],
			IsOpen:      c.IsOpen,
			IsSurvey:    c.IsSurvey,
		}
	}

	// tasksをchallengesと別に作っているのは、challengesはnot logged inなuserには見せてないから
	tasks := make([]string, len(challenges))
	for i, c := range challenges {
		tasks[i] = c.Name
	}

	// ----

	chalMap := make(map[uint]*Challenge)
	for _, c := range challenges {
		chalMap[c.ID] = c
	}

	teamSubmissions := make(map[uint][]*model.ValidSubmission)
	for _, t := range teams {
		teamSubmissions[t.ID] = make([]*model.ValidSubmission, 0)
	}
	for _, s := range submissions {
		teamSubmissions[s.TeamId] = append(teamSubmissions[s.TeamId], s)
	}

	// とりあえずエントリを作成する
	scoreFeed := make([]*ScoreFeedEntry, len(teams))
	for i := 0; i < len(teams); i++ {
		var score uint = 0
		taskStats := make(map[string]*TaskStat)
		var lastSubmission int64 = 0

		for _, s := range teamSubmissions[teams[i].ID] {
			c, exist := chalMap[s.ChallengeId]
			if !exist {
				continue //?
			}
			score += c.Score
			solvedAt := s.CreatedAt.Unix()
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

	// ----

	users, err := app.repo.ListAllUsers()
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
	}
	userSubmissions := make(map[uint][]*model.ValidSubmission)
	for _, u := range users {
		userSubmissions[u.ID] = make([]*model.ValidSubmission, 0)
	}
	for _, s := range submissions {
		userSubmissions[s.UserId] = append(userSubmissions[s.UserId], s)
	}

	// とりあえずエントリを作成する
	userScoreFeed := make([]*ScoreFeedEntry, len(users))
	for i := 0; i < len(users); i++ {
		var score uint = 0
		taskStats := make(map[string]*TaskStat)
		var lastSubmission int64 = 0

		for _, s := range userSubmissions[users[i].ID] {
			c, exist := chalMap[s.ChallengeId]
			if !exist {
				continue //?
			}
			score += c.Score
			solvedAt := s.CreatedAt.Unix()
			taskStats[c.Name] = &TaskStat{
				Score:    c.Score,
				SolvedAt: solvedAt,
			}
			if !c.IsSurvey && lastSubmission < solvedAt {
				lastSubmission = solvedAt
			}
		}

		userScoreFeed[i] = &ScoreFeedEntry{
			Pos:            0,
			Teamname:       users[i].Username,
			TeamID:         users[i].ID,
			Score:          int(score),
			TaskStats:      taskStats,
			LastSubmission: lastSubmission,
		}
	}

	// スコアと最終提出時刻でsort
	sort.Slice(userScoreFeed, func(i, j int) bool {
		if userScoreFeed[i].Score == userScoreFeed[j].Score {
			return userScoreFeed[i].LastSubmission < userScoreFeed[j].LastSubmission
		}
		return userScoreFeed[i].Score > userScoreFeed[j].Score
	})

	// Posの値を埋める
	for i := 0; i < len(userScoreFeed); i++ {
		userScoreFeed[i].Pos = i + 1
		if i != 0 && userScoreFeed[i].Score == userScoreFeed[i-1].Score && userScoreFeed[i].LastSubmission == userScoreFeed[i-1].LastSubmission {
			userScoreFeed[i].Pos = userScoreFeed[i-1].Pos
		}
	}

	return challenges, &Scoreboard{
			Tasks:     tasks,
			Standings: scoreFeed,
		}, &Scoreboard{
			Tasks:     tasks,
			Standings: userScoreFeed,
		}, nil
}
