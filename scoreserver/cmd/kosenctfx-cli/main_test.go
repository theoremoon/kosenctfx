package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestMakeDistFiles(t *testing.T) {
	tardata, err := makeDistfiles("./testdata/web/miniblog/distfiles", "miniblog_distfiles")
	if err != nil {
		t.Errorf("failed to run makeDistfiles: %+v\n", err.Error())
	}

	cmd := exec.Command("tar", "tfz", "-")
	cmd.Stdin = bytes.NewBuffer(tardata)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("the output is not tar formatted or cannot run tar: %+v\n", err)
	}

	expected := strings.TrimSpace(`
./miniblog_distfiles/miniblog/userdir/.keep
./miniblog_distfiles/miniblog/views/index.html
./miniblog_distfiles/miniblog/views/user.html
./miniblog_distfiles/miniblog/tmp/.keep
./miniblog_distfiles/miniblog/main.py
./miniblog_distfiles/miniblog/user_template/posts/00000000000000000000000000000000
./miniblog_distfiles/miniblog/user_template/titles/00000000000000000000000000000000
./miniblog_distfiles/miniblog/user_template/template
./miniblog_distfiles/miniblog/user_template/attachments/neko.png
`)

	if strings.TrimSpace(string(output)) != expected {
		t.Errorf("expected:\n%s\n\noutput:\n%s", expected, output)
	}
}

func TeetLoadTaskYaml(t *testing.T) {
	y, err := loadTaskYaml("./testdata/survey/survey/task.yml")
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	if !y.IsSurvey {
		t.Error("IsSurvey should be true\n")
	}
}
