package repository_test

import (
	"testing"

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

	return repository.New(db)
}

func TestMigrate(t *testing.T) {
	repo := newRepository()
	repo.Migrate()
}
