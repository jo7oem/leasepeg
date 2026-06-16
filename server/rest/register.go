package rest

import (
	"encoding/json"
	"net/http"
)

func NewRegisterHandler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /ping", GetState)
	m.HandleFunc("GET /echo", GetEcho)

	return m
}

func GetState(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func GetEcho(w http.ResponseWriter, r *http.Request) {
	rAddr := r.RemoteAddr

	resp := getEchoResponse{
		Address: rAddr,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type getEchoResponse struct {
	Address string `json:"address"`
}
