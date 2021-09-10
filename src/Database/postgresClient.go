package Database

import (
	"context"
	"errors"
	"fmt"
	"gourlshortener/src/Utilities"
	"reflect"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

func Validate() (*pgxpool.Config, error) {
	user := Utilities.ReadEnvironmentVariable("POSTGRES_USER", reflect.String, "")
	password := Utilities.ReadEnvironmentVariable("POSTGRES_PASSWORD", reflect.String, "")
	host := Utilities.ReadEnvironmentVariable("POSTGRES_HOST", reflect.String, "")
	port := Utilities.ReadEnvironmentVariable("POSTGRES_PORT", reflect.Int64, 0)
	dbName := Utilities.ReadEnvironmentVariable("POSTGRES_DATABASE", reflect.String, "")
	if user == "" || password == "" || host == "" || port.(int64) < 1 || dbName == "" {
		return nil, errors.New("PostgreSQL config parameter(s) invalid")
	}
	sslMode := Utilities.ReadEnvironmentVariable("POSTGRES_SSL_MODE", reflect.String, "")
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

func PerformDatabaseInsert(dbClient *pgxpool.Pool, link string, shortenedLink string) error {
	query := fmt.Sprintf(`SELECT 1 FROM urls WHERE original = '%s'`, link)
	rows, err := dbClient.Query(context.Background(), query)
	if err != nil {
		return err
	}
	if rows != nil {
		for rows.Next() {
			var number int = 0
			if rows.Scan(&number); number == 1 {
				return nil
			}
		}
	}
	query = fmt.Sprintf(`INSERT INTO urls (original, shortened) VALUES ('%s', '%s')`, link, shortenedLink)
	_, err = dbClient.Exec(context.Background(), query)
	if err != nil {
		return err
	}
	return nil
}

func PerformDatabaseSelect(dbClient *pgxpool.Pool, resolvable string) (string, error) {
	query := fmt.Sprintf(`SELECT original FROM urls WHERE shortened = '%s'`, resolvable)
	rows, err := dbClient.Query(context.Background(), query)
	if err != nil {
		return "", err
	}
	if rows != nil {
		var original string
		for rows.Next() {
			err := rows.Scan(&original)
			if err != nil {
				return "", err
			}
		}
		return original, nil
	}
	return "", nil
}
