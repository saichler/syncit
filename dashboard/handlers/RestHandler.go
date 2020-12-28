package handlers

import (
	"net/http"
)

type RestHandler interface {
	Endpoint() string
	Run(http.ResponseWriter, *http.Request)
	Method() string
}
