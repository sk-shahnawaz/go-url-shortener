package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"gourlshortener/src/Database"
	"gourlshortener/src/Handlers"
	"gourlshortener/src/Models"
	"gourlshortener/src/Utilities"
	_ "gourlshortener/src/docs"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title go-url-shortener API
// @version 1.0
// @description URL shortener & resolver service written in Golang using Echo web framework
// @contact.name Sk Shahnawaz-ul Haque
// @host localhost:3355
// @BasePath /api
// @schemes http https
func main() {
	var dbClient *pgxpool.Pool
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"error":   string(debug.Stack()),
				"message": "main crashed",
			}).Error()
		}
		if dbClient != nil {
			Database.Disconnect(dbClient)
		}
	}()

	godotenv.Load()
	app := echo.New()

	minimumLogLevel := Utilities.ReadEnvironmentVariable("LOG_MINIMUM_LEVEL", reflect.String, "ERROR")
	func(minimumLogLevel string) {
		logLevel, err := log.ParseLevel(minimumLogLevel)
		if err != nil {
			log.SetLevel(log.ErrorLevel)
		} else {
			log.SetLevel(logLevel)
		}
		log.SetOutput(os.Stdout)
	}(fmt.Sprintf("%v", minimumLogLevel))

	if useInMemoryDb := Utilities.ReadEnvironmentVariable("USE_IN_MEMORY_DB", reflect.String, "Y"); strings.ToUpper(useInMemoryDb.(string)) == "N" {
		dbClient, err = Database.Connect()
		if err != nil {
			log.Error("Failed to connect to database")
			panic("Failed to connect to database")
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
			extendedContext := Models.ExtendedContext{Context: c, Db: dbClient}
			return next(extendedContext)
		}
	})

	linkGeneratorRoute := apiGrouping.POST("/generate", Handlers.GenerateShortenedUrl)
	linkGeneratorRoute.Name = "Generator"

	linkResolverRoute := apiGrouping.GET("/resolve", Handlers.ResolveShortenedUrl)
	linkResolverRoute.Name = "Resolver"

	var portNumber int64 = 80
	if portNumberInt := Utilities.ReadEnvironmentVariable("CUSTOM_PORT_NUMBER", reflect.Int32, "80"); portNumberInt.(int64) > 0 {
		portNumber = portNumberInt.(int64)
	}

	app.GET("/swagger/*any", echoSwagger.WrapHandler)
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", portNumber)))
}
