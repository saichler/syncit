package handlers

import (
	"net/http"
)

type Ls struct {
}

func NewLs() *Ls {
	return &Ls{}
}

func (h *Ls) Endpoint() string {
	return "/dashboard/ls"
}

func (h *Ls) Run(w http.ResponseWriter, r *http.Request) {
	if !validateToken(r) {
		w.WriteHeader(401)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Hello!"))
}

func (h *Ls) Method() string {
	return "POST"
}
