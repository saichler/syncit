package dashboard

import (
	"github.com/gorilla/mux"
	"github.com/saichler/security"
	"github.com/saichler/syncit/dashboard/handlers"
	log "github.com/saichler/utils/golang"
	"net/http"
	"strconv"
)

const (
	CRT_FILE_PREFIX = "dashboard"
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
	ca, caKey, err := security.CreateCA(CRT_FILE_PREFIX, "syncit", "USA", "Santa Clara",
		"", "", "", "saichler@gmail.com", 10)
	if err != nil {
		log.Error(err)
	} else {
		err := security.CreateCrt(CRT_FILE_PREFIX, "syncit", "USA", "Santa Clara",
			"", "", "", "saichler@gmail.com", "127.0.0.1", "syncit", int64(port), 10, ca, caKey)
		log.Error(err)
	}
	return dashboard
}

func (dashboard *Dashboard) Start() {
	err := http.ListenAndServeTLS(":"+strconv.Itoa(dashboard.port), CRT_FILE_PREFIX+".crt", CRT_FILE_PREFIX+".crtKey", dashboard.router)
	if err != nil {
		panic(err.Error())
	}
}

func (dashBoard *Dashboard) RegisterHandler(handler handlers.RestHandler) {
	dashBoard.router.HandleFunc(handler.Endpoint(), handler.Run).Methods(handler.Method()).Schemes("https")
}
