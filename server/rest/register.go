package rest

import (
	"exp-myddns/logging"
	"exp-myddns/logic/rest"
	"log/slog"
	"net/http"
)

func NewRegisterHandler(logger slog.Handler) http.Handler {
	res := APIHandler{
		ServeMux: http.NewServeMux(),
		logger:   logging.NewLogger(logger, "restAPI"),
	}

	res.Handle("/", rest.NewHandler(logger, nil))

	return res
}

type APIHandler struct {
	*http.ServeMux

	logger *logging.Logger
}
