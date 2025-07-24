package core

import (
	"net/http"
)

type APIError int

const (
	ENotFound APIError = http.StatusNotFound
)

var apierror_message_mappings = map[APIError]string{
	ENotFound: "resource not found",
}

func (ae APIError) Error() string {
	msg, ok := apierror_message_mappings[ae]
	if ok {
		return msg
	}

	return "internal server error"
}
