package handlers

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/saichler/syncit/files"
	"github.com/saichler/syncit/model"
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

	file := &model.File{}
	err:=jsonpb.Unmarshal(r.Body, file)
	if err!=nil {
		fmt.Println(err)
	}
	file = files.Scan(file.NameA)

	resp, _ := model.PbMarshaler.MarshalToString(file)

	w.Write([]byte(resp))
}

func (h *Ls) Method() string {
	return "POST"
}
