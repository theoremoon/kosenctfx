package repository_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
)

func TestRegisterTeam(t *testing.T) {
	repo := newRepository()
	db := getDB()

	t.Run("登録できる", func(t *testing.T) {
		err := repo.RegisterTeam(&model.Team{
			Teamname:     "test",
			Email:        "test@example.com",
			PasswordHash: "somehash",
			CountryCode:  "",
		})
		if err != nil {
			t.Errorf("%+v\n", err)
		}

		var count int
		if err := db.QueryRow("select count(*) from teams").Scan(&count); err != nil {
			t.Fatalf("Query Error %+v\n", err)
		}
		if count == 0 {
			t.Errorf("this should be larger than 0")
		}
	})

	t.Run("teamnameがかぶっているときDuplicate Error", func(t *testing.T) {
		err := repo.RegisterTeam(&model.Team{
			Teamname:     "test",
			Email:        "test2@example.com",
			PasswordHash: "somehash",
			CountryCode:  "",
		})
		if err == nil {
			t.Errorf("this should be error %+v\n", err)
		}
	})

	t.Run("emailがかぶっているときDuplicate Error", func(t *testing.T) {
		err := repo.RegisterTeam(&model.Team{
			Teamname:     "test3",
			Email:        "test@example.com",
			PasswordHash: "somehash",
			CountryCode:  "",
		})
		if err == nil {
			t.Errorf("this should be error %+v\n", err)
		}
	})
}

func TestMakeTeamAdmin(t *testing.T) {
	repo := newRepository()
	team := &model.Team{
		Teamname:     gofakeit.Name(),
		Email:        gofakeit.Email(),
		PasswordHash: gofakeit.LetterN(16),
		CountryCode:  "",
	}
	err := repo.RegisterTeam(team)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	if team.IsAdmin {
		t.Errorf("this should be false")
	}

	err = repo.MakeTeamAdmin(team)
	if err != nil {
		t.Errorf("this should be nil %+v\n", err)
	}
	if !team.IsAdmin {
		t.Errorf("this should be true")
	}
}
