package dto

import (
	"github.com/go-ozzo/ozzo-validation/is"
	validator "github.com/go-ozzo/ozzo-validation/v4"
)

type Input struct {
	Url string `json:"url"`
}

func (input Input) Validate() error {
	return validator.ValidateStruct(&input,
		validator.Field(&input.Url, validator.Required.Error("'url' is required.")),
		validator.Field(&input.Url, is.URL.Error("Provide proper URL as value of 'url'")))
}

func Validate(queryString string) error {
	return validator.Validate(queryString, validator.Required.Error("Provide value of query string"))
}
