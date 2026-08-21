package utils

import (
	"net/http"

	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
)

func WrapErr(err error) *relaymodel.ErrorWithStatusCode {
	return &relaymodel.ErrorWithStatusCode{
		StatusCode: http.StatusInternalServerError,
		Error: relaymodel.Error{
			Message: err.Error(),
		},
	}
}
