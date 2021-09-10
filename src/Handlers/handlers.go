package Handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gourlshortener/src/Database"
	"gourlshortener/src/Models"
	"gourlshortener/src/Models/DTOs"
	"gourlshortener/src/Utilities"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var store map[string]string = make(map[string]string)

// Shortened URL generator API godoc
// @Summary Shortened URL generator
// @Description Generates shortened URL
// @Tags Generate
// @accept json
// @Param input body DTOs.Input true "Input"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /generate [post]
func GenerateShortenedUrl(context echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			err := errors.New("GenerateShortenedUrl crashed")
			log.WithFields(log.Fields{
				"stacktrace": string(debug.Stack()),
				"err":        err,
			}).Error()
		}
	}()
	log.WithFields(log.Fields{
		"path":   "/api/generate",
		"method": http.MethodPost,
		"name":   "Generate Shortened URL",
	}).Info()
	if context.Request().Body == nil {
		log.WithFields(log.Fields{
			"status": http.StatusBadRequest,
		}).Warn()
		return context.String(http.StatusBadRequest, "Empty body.")
	}
	if log.GetLevel() == log.DebugLevel {
		serializedByteContent, err := json.Marshal(context.Request().Body)
		if err != nil {
			log.WithFields(log.Fields{
				"status":     http.StatusInternalServerError,
				"stackTrace": string(debug.Stack()),
			}).Warn()
		}
		log.WithFields(log.Fields{
			"body": string(serializedByteContent),
		}).Debug()
	}
	input := new(DTOs.Input)
	if err := (&echo.DefaultBinder{}).BindBody(context, &input); input == nil || err != nil {
		if err != nil {
			log.WithFields(log.Fields{
				"status":     http.StatusInternalServerError,
				"stackTrace": string(err.Error()),
			}).Error()
			return context.String(http.StatusInternalServerError, "Error occurred while binding request body.")
		} else if input == nil {
			log.WithFields(log.Fields{
				"status": http.StatusBadRequest,
			}).Error()
			return context.String(http.StatusBadRequest, "Bad request received.")
		}
	}
	if err := input.Validate(); err != nil {
		log.WithFields(log.Fields{
			"status":     http.StatusInternalServerError,
			"stackTrace": string(err.Error()),
		}).Error()
		return context.String(http.StatusInternalServerError, err.Error())
	}
	shortnedLink, err := Utilities.GenerateShortLink(input.Url)
	if err != nil {
		log.WithFields(log.Fields{
			"status":     http.StatusInternalServerError,
			"stackTrace": string(err.Error()),
		}).Error()
		return context.String(http.StatusInternalServerError, "Error occurred shortening URL.")
	}
	go func() {
		if useInMemoryDb := Utilities.ReadEnvironmentVariable("USE_IN_MEMORY_DB", reflect.String, "Y"); strings.ToUpper(useInMemoryDb.(string)) == "N" {
			err = Database.PerformDatabaseInsert(context.(Models.ExtendedContext).Db, input.Url, shortnedLink)
		} else {
			if _, present := store[input.Url]; !present {
				store[input.Url] = shortnedLink
			}
		}
	}()
	log.WithFields(log.Fields{
		"status":       http.StatusOK,
		"originalUrl":  input.Url,
		"shortenedUrl": shortnedLink,
	}).Debug()
	return context.String(http.StatusOK, fmt.Sprint("http://", context.Request().Host, "/api/resolve", "?q=", shortnedLink))
}

// Shortened URL resolver API godoc
// @Summary Shortened URL resolver
// @Description Resolves the shortened URL and redirects to resolved URL
// @Tags Resolve
// @param q query string true "q is mandatory"
// @Produce  html
// @Success 308 {object} string
// @Failure 500 {object} string
// @Router /resolve [get]
func ResolveShortenedUrl(context echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			err := errors.New("ResolveShortenedUrl crashed")
			log.WithFields(log.Fields{
				"stacktrace": string(debug.Stack()),
				"err":        err,
			}).Error()
		}
	}()
	log.WithFields(log.Fields{
		"path":   "/api/resolve",
		"method": http.MethodPost,
		"name":   "Resolve Shortened URL",
	}).Info()
	queryParameterValue := context.QueryParam("q")
	if err := DTOs.Validate(queryParameterValue); err != nil {
		log.WithFields(log.Fields{
			"status":  http.StatusBadRequest,
			"message": err,
		}).Error()
		return context.String(http.StatusInternalServerError, "Query string missing.")
	}
	resolvedLink, err := resolveShortenedLink(context.(Models.ExtendedContext).Db, queryParameterValue)
	if err != nil {
		log.WithFields(log.Fields{
			"status":     http.StatusInternalServerError,
			"stackTrace": string(err.Error()),
		}).Error()
		return context.String(http.StatusInternalServerError, "Error occurred resolving URL.")
	}
	return context.Redirect(http.StatusPermanentRedirect, resolvedLink)
}

func resolveShortenedLink(dbClient *pgxpool.Pool, resolvable string) (string, error) {
	if useInMemoryDb := Utilities.ReadEnvironmentVariable("USE_IN_MEMORY_DB", reflect.String, "Y"); strings.ToUpper(useInMemoryDb.(string)) == "N" {
		return Database.PerformDatabaseSelect(dbClient, resolvable)
	} else {
		for link, shortenedLink := range store {
			if shortenedLink == resolvable {
				return link, nil
			}
		}
		return "", errors.New("no entry found")
	}
}
