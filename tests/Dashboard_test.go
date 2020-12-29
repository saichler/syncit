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
	resp, e := execute(h, payload, "", "", t)
	if e != nil {
		return
	}

	jsonpb.UnmarshalString(resp, userpass)

	if userpass.Token == "" {
		fmt.Println("token is blank")
		t.Fail()
		return
	}

	filename := "/home/saichler/syncit"
	file := &model.File{}
	file.NameA = filename
	payload, _ = model.PbMarshaler.MarshalToString(file)

	h2 := handlers.NewLs()
	resp, e = execute(h2, payload, userpass.Token, "", t)
	if e != nil {
		return
	}
	jsonpb.UnmarshalString(resp, file)
	files.Print(file, 2, false, true)

	fmt.Println("Get with parameters:")

	resp, e = execute(h2, "", userpass.Token, "file=/home/saichler/syncit;dept=4;incFile=true;incLessBlock=true", t)

	fmt.Println(resp)

	//time.Sleep(time.Second * 60)
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

func execute(h handlers.RestHandler, payload, token, arg string, t *testing.T) (string, error) {
	hc := &http.Client{}
	url := "https://127.0.0.1:10101" + h.Endpoint()
	if arg != "" {
		url += "?" + arg
	}
	request, e := createRequest(h.Method(), url, []byte(payload), token)
	if e != nil {
		t.Fail()
		fmt.Println(e)
		return "", e
	}
	resp, e := hc.Do(request)
	if e != nil {
		t.Fail()
		fmt.Println(e)
		return "", e
	}
	jsn, _ := ioutil.ReadAll(resp.Body)
	jsonStr := string(jsn)
	return jsonStr, nil
}
