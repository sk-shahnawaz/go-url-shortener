package main

import (
	"fmt"
	"gourlshortener/db"
	"gourlshortener/utilities"
	"gourlshortener/web-api/handlers"
	"gourlshortener/web-api/models"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	var dbClient *pgxpool.Pool
	var err error
	defer func() {
		if dbClient != nil {
			db.Disconnect(dbClient)
		}
	}()

	godotenv.Load()
	app := echo.New()

	minimumLogLevel := utilities.ReadEnvironmentVariable("LOG_MINIMUM_LEVEL", reflect.String, "ERROR")
	func(minimumLogLevel string) {
		logLevel, err := log.ParseLevel(minimumLogLevel)
		if err != nil {
			log.SetLevel(log.ErrorLevel)
		} else {
			log.SetLevel(logLevel)
		}
		log.SetOutput(os.Stdout)
	}(fmt.Sprintf("%v", minimumLogLevel))

	if useInMemoryDb := utilities.ReadEnvironmentVariable("USE_IN_MEMORY_DB", reflect.String, "Y"); strings.ToUpper(useInMemoryDb.(string)) == "N" {
		dbClient, err = db.Connect()
		if err != nil {
			log.Error("Failed to connect to database")
			os.Exit(1)
		}
	}

	app.Use(middleware.Recover())
	app.Use(middleware.RemoveTrailingSlash())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	app.GET("/", func(context echo.Context) error {
		log.WithFields(log.Fields{
			"path":   "/",
			"method": http.MethodGet,
			"name":   "Health",
			"status": http.StatusOK,
		}).Debug()
		return context.String(http.StatusOK, "Server up & running.")
	}).Name = "Health"

	apiGrouping := app.Group("/api")

	apiGrouping.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			extendedContext := models.ExtendedContext{Context: c, Db: dbClient}
			return next(extendedContext)
		}
	})

	linkGenerationRoute_Generator := apiGrouping.POST("/generate", handlers.GenerateShortenedUrl)
	linkGenerationRoute_Generator.Name = "Generator"

	linkGenerationRoute_Resolver := apiGrouping.GET("/resolve", handlers.ResolveShortenedUrl)
	linkGenerationRoute_Resolver.Name = "Resolver"

	var portNumber int64 = 80
	if portNumberInt := utilities.ReadEnvironmentVariable("CUSTOM_PORT_NUMBER", reflect.Int32, "80"); portNumberInt.(int64) > 0 {
		portNumber = portNumberInt.(int64)
	}

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", portNumber)))
}
