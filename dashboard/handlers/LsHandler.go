package handlers

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/saichler/syncit/files"
	"github.com/saichler/syncit/model"
	"io/ioutil"
	"net/http"
	"strconv"
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
	data, e := ioutil.ReadAll(r.Body)
	jsonPayload := true
	if e == nil && data != nil && len(data) > 0 {
		err := jsonpb.UnmarshalString(string(data), file)
		if err != nil {
			w.Write([]byte(err.Error()))
			fmt.Println(err)
		}
	} else if data != nil && len(data) == 0 {
		f := r.URL.Query().Get("file")
		if f == "" {
			w.Write([]byte("No file was specified in the query"))
			return
		}
		file.NameA = f
		jsonPayload = false
	}
	file = files.Scan(file.NameA)

	if jsonPayload {
		resp, _ := model.PbMarshaler.MarshalToString(file)
		w.Write([]byte(resp))
	} else {
		dept := 2
		incFile := true
		incLessThanBlock := true
		dpt := r.URL.Query().Get("dept")
		if dpt != "" {
			dept, _ = strconv.Atoi(dpt)
		}
		incF := r.URL.Query().Get("incFile")
		if incF != "" {
			if incF == "true" {
				incFile = true
			} else {
				incFile = false
			}
		}
		incB := r.URL.Query().Get("incLessBlock")
		if incB != "" {
			if incB == "true" {
				incLessThanBlock = true
			} else {
				incLessThanBlock = false
			}
		}
		buff := files.Print(file, dept, incFile, incLessThanBlock)
		w.Write(buff.Bytes())
	}
}

func (h *Ls) Method() string {
	return "GET"
}
