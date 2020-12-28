package handlers

import (
	"crypto/md5"
	"encoding/hex"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/jsonpb"
	"github.com/saichler/security"
	"github.com/saichler/syncit/model"
	log "github.com/saichler/utils/golang"
	"net/http"
	"time"
)

var secret = "sync-it"

type Login struct {
}

func NewLogin() *Login {
	return &Login{}
}

func (h *Login) Endpoint() string {
	return "/dashboard/login"
}

type Claim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (h *Login) Run(w http.ResponseWriter, r *http.Request) {
	model.InitSt()

	userpass := &model.UserPass{}
	err := jsonpb.Unmarshal(r.Body, userpass)

	st := security.InitSecureStore(model.IO_FILENAME)
	passHash, err := st.Get("/users/" + userpass.Username)
	if err != nil || passHash == "" {
		w.WriteHeader(401)
		return
	}

	hash := md5.New()
	md5Hash := hex.EncodeToString(hash.Sum([]byte(userpass.Password)))

	if md5Hash != passHash {
		w.WriteHeader(401)
		return
	}

	claim := &Claim{}
	claim.Username = userpass.Username
	claim.ExpiresAt = time.Now().Add(time.Minute * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(model.Secret))
	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(200)
	userpass.Username = ""
	userpass.Password = ""
	userpass.Token = tokenStr
	response, _ := model.PbMarshaler.MarshalToString(userpass)
	_, e := w.Write([]byte(response))
	if e != nil {
		log.Error(err)
	}
}

func (h *Login) Method() string {
	return "POST"
}

func validateToken(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		return false
	}
	claim := &Claim{}
	tkn, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(model.Secret), nil
	})
	if err != nil {
		return false
	}
	if !tkn.Valid {
		return false
	}
	return true
}
