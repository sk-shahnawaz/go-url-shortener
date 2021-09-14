package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
	log "github.com/sirupsen/logrus"
)

func init() {
	up := func(db orm.DB) error {
		log.Info("Creating tables..")
		_, err := db.Exec(`
			CREATE TABLE urls (original VARCHAR(512), shortened VARCHAR(50));
		`)
		if err != nil {
			log.Error("Error occurred while creating tables!")
		}
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE urls")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210914143447_create-schema", up, down, opts)
}
