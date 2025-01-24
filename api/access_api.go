package api

import (
	"github.com/gobwas/ws"
	"log"
	"net/http"
)

type AccessAPI struct {
	WSManager *WSManager
}

func NewAccessAPI(buffer chan []byte) *AccessAPI {
	return &AccessAPI{
		WSManager: NewWSManager(buffer),
	}
}

func (a *AccessAPI) Start() {
	http.HandleFunc("/ws", http.HandlerFunc(a.handleWS))

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func (a *AccessAPI) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println(err)
		return
	}

	a.WSManager.RegisterWSConnection(conn)
}
