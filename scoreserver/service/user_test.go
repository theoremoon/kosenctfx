package service

import "testing"

func Test_RegsiterUserWithTeam(t *testing.T) {
	app := newTestApp(t)
	defer app.Close()

	t := model.Team{
		Teamname: "team1",
		Token:    "some_unique_token",
	}
	err := app.RegisterTeam(&t)
	assert.NoError(err)
}
