package repository_test

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/theoremoon/kosenctfx/scoreserver/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	testDB = "testuser:testpassword@tcp(db-test:3306)/testtable"
)

func newRepository() repository.Repository {
	db, err := gorm.Open(mysql.Open(testDB), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	repo := repository.New(db)
	repo.Migrate()
	return repo
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", testDB)
	if err != nil {
		panic(err)
	}
	return db
}
