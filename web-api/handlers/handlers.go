package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gourlshortener/utilities"
	"gourlshortener/web-api/dto"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

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
		return context.JSON(http.StatusBadRequest, "Empty body.")
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
	input := new(dto.Input)
	if err := (&echo.DefaultBinder{}).BindBody(context, &input); input == nil || err != nil {
		if err != nil {
			log.WithFields(log.Fields{
				"status":     http.StatusInternalServerError,
				"stackTrace": string(err.Error()),
			}).Error()
			return context.JSON(http.StatusInternalServerError, "Error occurred while binding request body.")
		} else if input == nil {
			log.WithFields(log.Fields{
				"status": http.StatusBadRequest,
			}).Error()
			return context.JSON(http.StatusBadRequest, "Bad request received.")
		}
	}
	shortnedLink, err := utilities.GenerateShortLink(input.Url)
	if err != nil {
		log.WithFields(log.Fields{
			"status":     http.StatusInternalServerError,
			"stackTrace": string(err.Error()),
		}).Error()
		return context.JSON(http.StatusInternalServerError, "Error occurred shortening URL.")
	}
	log.WithFields(log.Fields{
		"status":       http.StatusOK,
		"originalUrl":  input.Url,
		"shortenedUrl": shortnedLink,
	}).Debug()
	return context.String(http.StatusOK, fmt.Sprint("http://", context.Request().Host, "/api/resolver", "?q=", shortnedLink))
}

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
	if queryParameterValue == "" {
		log.WithFields(log.Fields{
			"status":  http.StatusBadRequest,
			"message": "Query parameter missing",
		}).Debug()
	}
	resolvedLink, err := utilities.ResolveShortenedLink(queryParameterValue)
	if err != nil {
		log.WithFields(log.Fields{
			"status":     http.StatusInternalServerError,
			"stackTrace": string(err.Error()),
		}).Error()
		return context.JSON(http.StatusInternalServerError, "Error occurred resolving URL.")
	}
	return context.Redirect(http.StatusPermanentRedirect, resolvedLink)
}
