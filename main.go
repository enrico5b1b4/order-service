package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/enrico5b1b4/order-service/app"
	"github.com/enrico5b1b4/order-service/env"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	port := env.MustGetEnv("PORT")
	postgresDSN := env.MustGetEnv("POSTGRES_DSN")
	orderProcessServiceURL := env.MustGetEnv("ORDER_PROCESS_SERVICE_URL")
	completeOrderCallbackURL := env.MustGetEnv("COMPLETE_ORDER_CALLBACK_URL")

	db, err := sqlx.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalln(err)
	}

	errWait := waitForDB(db)
	if errWait != nil {
		log.Fatalln(errWait)
	}

	// migrations
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalln(err)
	}
	_ = m.Up()

	orderServiceApp := app.New(db, orderProcessServiceURL, completeOrderCallbackURL)
	orderServiceApp.Logger.Fatal(orderServiceApp.Start(fmt.Sprintf(":%s", port)))
}

func waitForDB(db *sqlx.DB) error {
	var err error
	for i := 1; i < 5; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		log.Println("DB not ready. Waiting 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	return err
}
