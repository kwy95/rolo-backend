package api

import (
	"github.com/gobwas/ws"
	"log"
	"net/http"
	"rolo/backend/racer"
)

type AccessAPI struct {
	WSManager *WSManager
	input     chan []byte
}

func NewAccessAPI(in chan []byte, out chan []byte) *AccessAPI {
	return &AccessAPI{
		WSManager: NewWSManager(out),
		input:     in,
	}
}

func (a *AccessAPI) Start() {
	http.HandleFunc("/ws", http.HandlerFunc(a.handleWS))
	http.HandleFunc("/start", http.HandlerFunc(a.handleStart))

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

func (a *AccessAPI) handleStart(w http.ResponseWriter, r *http.Request) {
	race := racer.NewRace(a.input, a.WSManager.buffer)
	race.Start()
}
