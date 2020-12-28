package model

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/golang/protobuf/jsonpb"
	"github.com/saichler/security"
	"os"
)

var Secret = ""

const (
	IO_FILENAME = "./syncit.io"
	MYK         = "/my/k"
	MYS         = "/my/s"
	MYP         = "/my/p"
)

var PbMarshaler = &jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: false,
	OrigName:     true,
}

func InitSt() {
	_, err := os.Stat(IO_FILENAME)
	st := security.InitSecureStore(IO_FILENAME)
	if err != nil {
		st.Put(MYK, security.GenerateAES256Key())
		st.Put(MYS, "sync-it")
		st.Put(MYP, "45454")
		hash := md5.New()
		md5Hash := hex.EncodeToString(hash.Sum([]byte("world")))
		st.Put("/users/hello", md5Hash)
	}
	s, _ := st.Get(MYS)
	Secret = s
}
