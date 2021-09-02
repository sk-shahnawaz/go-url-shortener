package dto

type Input struct {
	Url string `json:"url" validate:"required"`
}
