package dashboard

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/gorilla/mux"
	"github.com/saichler/syncit/dashboard/handlers"
	log "github.com/saichler/utils/golang"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Dashboard struct {
	router *mux.Router
	port   int
}

func NewDashboard(port int) *Dashboard {
	dashboard := &Dashboard{}
	dashboard.router = mux.NewRouter()
	dashboard.port = port
	dashboard.RegisterHandler(handlers.NewLogin())
	dashboard.RegisterHandler(handlers.NewLs())
	err := dashboard.generate()
	if err != nil {
		log.Error(err)
	}
	return dashboard
}

func (dashboard *Dashboard) Start() {
	err := http.ListenAndServeTLS(":"+strconv.Itoa(dashboard.port), "/tmp/server.crt", "/tmp/server.key", dashboard.router)
	if err != nil {
		panic(err.Error())
	}
}

func (dashBoard *Dashboard) RegisterHandler(handler handlers.RestHandler) {
	dashBoard.router.HandleFunc(handler.Endpoint(), handler.Run).Methods(handler.Method()).Schemes("https")
}

func (dashBoard *Dashboard) generate() error {
	_, e := os.Stat("/tmp/server.crt")
	if e != nil {
		ca := &x509.Certificate{
			SerialNumber: big.NewInt(2019),
			Subject: pkix.Name{
				Organization:  []string{"Sync-it"},
				Country:       []string{"US"},
				Province:      []string{"Santa Clara"},
				Locality:      []string{"San Jose"},
				StreetAddress: []string{"1993 Curtner Ave"},
				PostalCode:    []string{"95124"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
			EmailAddresses:        []string{"saichler@gmail.com"},
		}

		caKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return err
		}

		caData, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
		if err != nil {
			return err
		}

		caPEM := &bytes.Buffer{}
		pem.Encode(caPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caData,
		})

		ioutil.WriteFile("/tmp/server.crt", caPEM.Bytes(), 0777)

		caKeyPEM := &bytes.Buffer{}
		pem.Encode(caKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(caKey),
		})
		ioutil.WriteFile("/tmp/server.key", caKeyPEM.Bytes(), 0777)
	}
	return nil
}
