package api

import (
	"github.com/gobwas/ws"
	"log"
	"net/http"
	"rolo/backend/racer"
)

type AccessAPI struct {
	WSManager *WSManager
	Race      *racer.Race
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
	http.HandleFunc("/stop", http.HandlerFunc(a.handleStop))

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
	if a.Race != nil {
		a.Race.EndRace()
		a.Race = nil
	}
	log.Println("Race started")
	a.Race = racer.NewRace(a.input, a.WSManager.buffer)
	a.Race.Start()
}

func (a *AccessAPI) handleStop(w http.ResponseWriter, r *http.Request) {
	if a.Race == nil {
		log.Println("No race runnning")
		return
	}
	log.Println("Race stopped")
	a.Race.EndRace()
}
