package handlers

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInternalServerError = errors.New("internal server error")
)