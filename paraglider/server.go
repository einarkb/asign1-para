package paragliding

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Server struct {
	db          *Database
	Mgrticker   *MgrTicker
	MgrWebhooks *WebHookMgr
	MgrTrack    *TrackMgr
}

func (server *Server) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func (server *Server) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT not set")
		return
	}

	http.HandleFunc("/", server.Hello)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
