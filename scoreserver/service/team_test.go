package service_test

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestRegister(t *testing.T) {
	app := newApp(t)

	_, err := app.RegisterTeam("team", "password", "team@example.com", "")
	assert.NilError(t, err)
}
