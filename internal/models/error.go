package models

type ErrorResponce struct {
	Errors string `json:"errors"`
}

type ServiceError struct {
	TextError string
	Code      int
}
