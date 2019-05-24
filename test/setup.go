package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var DSN string
var testDB *sqlx.DB
var err error

// Get correct path for migrations
var _, b, _, _ = runtime.Caller(0)
var basepath = filepath.Dir(b)

func init() {
	DSN = os.Getenv("POSTGRES_TEST_DSN")

	if DSN != "" && testDB == nil {
		testDB, err = sqlx.Connect("postgres", DSN)
		if err != nil {
			log.Fatalln(err)
		}
		_ = waitForDB(testDB)

		// migration
		driver, err := postgres.WithInstance(testDB.DB, &postgres.Config{})
		migrationsDirectory := fmt.Sprintf("file://%s/../migrations", basepath)
		m, err := migrate.NewWithDatabaseInstance(migrationsDirectory, "postgres", driver)
		if err != nil {
			panic(err)
		}

		_ = m.Up()
	}
}

func DBSetup(setup func(db *sqlx.DB)) *sqlx.DB {
	setup(testDB)
	return testDB
}

func DBConnect() *sqlx.DB {
	_ = waitForDB(testDB)

	return testDB
}

func CheckSkipTest(t *testing.T) {
	checkDSN := os.Getenv("POSTGRES_TEST_DSN")
	if checkDSN == "" {
		t.Skip()
		return
	}
}

func waitForDB(db *sqlx.DB) error {
	var err error
	count := 5
	for i := 1; i < count; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		log.Println("DB not ready. Waiting 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	return err
}
