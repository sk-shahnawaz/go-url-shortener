package db

import (
	"context"
	"errors"
	"fmt"
	"gourlshortener/utilities"
	"reflect"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

func Validate() (*pgxpool.Config, error) {
	user := utilities.ReadEnvironmentVariable("POSTGRES_USER", reflect.String, "")
	password := utilities.ReadEnvironmentVariable("POSTGRES_PASSWORD", reflect.String, "")
	host := utilities.ReadEnvironmentVariable("POSTGRES_HOST", reflect.String, "")
	port := utilities.ReadEnvironmentVariable("POSTGRES_PORT", reflect.Int64, 0)
	dbName := utilities.ReadEnvironmentVariable("POSTGRES_DATABASE", reflect.String, "")
	if user == "" || password == "" || host == "" || port.(int64) < 1 || dbName == "" {
		return nil, errors.New("PostgreSQL config parameter(s) invalid")
	}
	sslMode := utilities.ReadEnvironmentVariable("POSTGRES_SSL_MODE", reflect.String, "")
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, password, host, port, dbName)
	if sslMode != "" {
		connectionString = fmt.Sprintf("%v?sslmode=%v", connectionString, sslMode)
	}
	config, err := pgxpool.ParseConfig(connectionString)
	if config != nil && err == nil {
		log.WithFields(log.Fields{
			"connectionString": connectionString,
		}).Debug("Parsed the db configuration successfully.")
		return config, nil
	}
	return nil, err
}

func Connect() (*pgxpool.Pool, error) {
	config, err := Validate()
	if err != nil {
		log.WithFields(log.Fields{
			"connectionString": config.ConnString(),
			"err":              err.Error(),
		}).Error("Error while trying to valiate configs")
		return nil, err
	}
	log.WithFields(log.Fields{
		"connectionString": config.ConnString(),
	}).Debug("Trying to open DB connection")

	dbClient, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}
	log.Debug("Postgres: Created the connection pool successfully.")
	return dbClient, err
}

func Disconnect(dbClient *pgxpool.Pool) {
	dbClient.Close()
}
