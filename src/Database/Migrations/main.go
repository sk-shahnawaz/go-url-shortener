package main

import (
	"fmt"
	"gourlshortener/src/Utilities"
	"log"
	"os"
	"reflect"

	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func main() {
	godotenv.Load("./../../.env")

	user := Utilities.ReadEnvironmentVariable("POSTGRES_USER", reflect.String, "")
	password := Utilities.ReadEnvironmentVariable("POSTGRES_PASSWORD", reflect.String, "")
	host := Utilities.ReadEnvironmentVariable("POSTGRES_HOST", reflect.String, "")
	port := Utilities.ReadEnvironmentVariable("POSTGRES_PORT", reflect.Int64, 0)
	dbName := Utilities.ReadEnvironmentVariable("POSTGRES_DATABASE", reflect.String, "")

	if user == "" || password == "" || host == "" || port.(int64) < 1 || dbName == "" {
		panic("PostgreSQL config parameter(s) invalid")
	}

	fmt.Println(user, host, port, dbName)

	portNumber := fmt.Sprintf("%d", port)
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", host.(string), portNumber),
		User:     user.(string),
		Database: dbName.(string),
		Password: password.(string),
	})

	const pathToMigrationScripts = "Database/Migrations"
	err := migrations.Run(db, pathToMigrationScripts, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
