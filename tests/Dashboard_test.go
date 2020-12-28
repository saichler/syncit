package tests

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/saichler/syncit/dashboard"
	"github.com/saichler/syncit/dashboard/handlers"
	"github.com/saichler/syncit/files"
	"github.com/saichler/syncit/model"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestDashboard(t *testing.T) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	d := dashboard.NewDashboard(10101)
	go d.Start()
	time.Sleep(time.Second)
	hc := &http.Client{}

	userpass := &model.UserPass{}
	userpass.Username = "hello"
	userpass.Password = "world"

	payload, e := model.PbMarshaler.MarshalToString(userpass)
	if e != nil {
		t.Fail()
		fmt.Println(e)
		return
	}

	h := &handlers.Login{}

	request, e := createRequest("POST", "https://127.0.0.1:10101"+h.Endpoint(), []byte(payload), "")
	if e != nil {
		t.Fail()
		fmt.Println(e)
		return
	}
	resp, e := hc.Do(request)
	if e != nil {
		t.Fail()
		fmt.Println(e)
		return
	}
	jsn, _ := ioutil.ReadAll(resp.Body)
	jsonStr := string(jsn)
	jsonpb.UnmarshalString(jsonStr, userpass)
	if userpass.Token == "" {
		t.Fail()
		return
	}

	h2 := handlers.NewLs()

	filename := "/home/saichler/syncit"
	file := &model.File{}
	file.NameA = filename
	payload, _ = model.PbMarshaler.MarshalToString(file)
	request, e = createRequest("POST", "https://127.0.0.1:10101"+h2.Endpoint(), []byte(payload), userpass.Token)
	if e != nil {
		fmt.Println(e)
		t.Fail()
		return
	}
	resp, e = hc.Do(request)
	if e != nil {
		fmt.Println(e)
		t.Fail()
		return
	}
	jsn, _ = ioutil.ReadAll(resp.Body)
	jsonStr = string(jsn)
	jsonpb.UnmarshalString(jsonStr, file)
	files.Print(file, 2, false, true)
}

func createRequest(rest, url string, body []byte, token string) (*http.Request, error) {
	request, err := http.NewRequest(rest, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", token)
	//request.Header.Add("content-type", "application/json; charset=UTF-8")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	return request, nil
}
