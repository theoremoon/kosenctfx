package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/theoremoon/kosenctfx/scoreserver/config"
	"github.com/theoremoon/kosenctfx/scoreserver/model"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/seeder"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

func run() error {
	rand.Seed(time.Now().UnixNano())

	size := flag.Int("size", 0, "The size of competition")
	flag.Usage = func() {
		fmt.Printf("%s\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if size == nil {
		flag.Usage()
		return nil
	}

	conf, err := config.Load()
	if err != nil {
		return err
	}
	db, err := gorm.Open(mysql.Open(conf.Dbdsn), &gorm.Config{})
	if err != nil {
		return err
	}
	rawdb, err := db.DB()
	if err != nil {
		return err
	}
	defer rawdb.Close()

	repo := repository.New(db)
	repo.Migrate()

	app := service.New(repo, nil)
	s := seeder.New(app)

	numofTeams := rand.Intn(*size)
	teams := make([]*model.Team, numofTeams)
	for i := 0; i < numofTeams; i++ {
		t, err := s.Team()
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		teams[i] = t
	}

	flagFormat := "flag{%s}"
	numofChallenges := rand.Intn(*size)
	challenges := make([]*model.Challenge, numofChallenges)
	for i := 0; i < numofChallenges; i++ {
		c, err := s.Challenge(flagFormat)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		challenges[i] = c
	}

	numofSubmissions := rand.Intn(numofChallenges * numofTeams)
	for i := 0; i < numofSubmissions; i++ {
		is_correct := rand.Intn(2)
		tid := rand.Intn(numofTeams)
		cid := rand.Intn(numofChallenges)
		if is_correct == 0 { // correct
			if _, err := s.Submission(teams[tid], challenges[cid]); err != nil {
				return xerrors.Errorf(": %w", err)
			}
		} else { // wrong
			if _, err := s.Submission(teams[tid], nil); err != nil {
				return xerrors.Errorf(": %w", err)
			}
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
