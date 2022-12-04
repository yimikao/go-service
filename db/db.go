package db

import (
	"log"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

func New() (*pg.DB, error) {
	var (
		opts       *pg.Options
		err        error
		dbConn     *pg.DB
		collection *migrations.Collection
	)

	if os.Getenv("ENV") == "prod" {
		opts, err = pg.ParseURL(os.Getenv("DB_URL"))
		if err != nil {
			return nil, err
		}
	} else {
		opts = &pg.Options{
			// Addr:     "db:5432",
			Addr:     "localhost:5430",
			User:     "postgres",
			Password: "postgres",
			Database: "postgrestut",
		}
	}

	dbConn = pg.Connect(opts)

	collection = migrations.NewCollection()
	err = collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return nil, err
	}

	_, _, err = collection.Run(dbConn, "reset")
	if err != nil {
		return nil, err
	}

	_, _, err = collection.Run(dbConn, "down")
	if err != nil {
		return nil, err
	}

	newV, oldV, err := collection.Run(dbConn, "up")
	if err != nil {
		return nil, err
	}

	if newV != oldV {
		log.Printf("changes made in migrations. upgraded from %d, to %d", oldV, newV)
	} else {
		log.Printf("no changes made to db migrations. current version: %d", oldV)
	}

	return dbConn, nil

}
