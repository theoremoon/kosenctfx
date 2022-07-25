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
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/seeder"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

func run() error {
	rand.Seed(time.Now().UnixNano())

	size := flag.Int("size", 0, "The size of competition")
	team := flag.Bool("team", false, "Seeding teams")
	challenge := flag.Bool("challenge", false, "Seeding challenges")
	submission := flag.Bool("submission", false, "Seeding submissios")
	all := flag.Bool("all", false, "Seed all values")
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

	app := service.New(db, nil)
	s := seeder.New(app)

	if *all || *team {
		n := rand.Intn(*size)
		for i := 0; i < n; i++ {
			_, err := s.Team()
			if err != nil {
				return xerrors.Errorf(": %w", err)
			}
		}
	}

	if *all || *challenge {
		flagFormat := "flag{%s}"
		n := rand.Intn(*size)
		for i := 0; i < n; i++ {
			_, err := s.Challenge(flagFormat)
			if err != nil {
				return xerrors.Errorf(": %w", err)
			}
		}
	}

	if *all || *submission {
		now := time.Now().Unix()

		n := rand.Intn(*size)
		teams, err := app.ListTeams()
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		challenges, err := app.ListAllRawChallenges()
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		for i := 0; i < n; i++ {
			is_correct := rand.Intn(2)
			tid := rand.Intn(len(teams))
			cid := rand.Intn(len(challenges))
			if is_correct != 0 { // correct
				if _, err := s.Submission(teams[tid], challenges[cid], now); err != nil {
					return xerrors.Errorf(": %w", err)
				}
			} else { // wrong
				if _, err := s.Submission(teams[tid], nil, now); err != nil {
					return xerrors.Errorf(": %w", err)
				}
			}

			now += int64(rand.Intn(3600))
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
