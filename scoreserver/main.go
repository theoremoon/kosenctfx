package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/theoremoon/kosenctfx/scoreserver/config"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"github.com/theoremoon/kosenctfx/scoreserver/server"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

func run() error {
	conf, err := config.Load()
	if err != nil {
		return err
	}
	db, err := gorm.Open("mysql", conf.Dbdsn)
	if err != nil {
		return err
	}
	defer db.Close()
	db.BlockGlobalUpdate(true)

	repo := repository.New(db)
	repo.Migrate()
	app := service.New(repo)
	srv := server.New(app)
	return srv.Start(conf.Addr)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
