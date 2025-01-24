package api

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"sync"
)

type WSManager struct {
	live   *liveSocket
	buffer chan []byte
	lock   sync.Mutex
}

type liveSocket struct {
	conn net.Conn
}

func NewWSManager(buffer chan []byte) *WSManager {
	m := &WSManager{
		live:   nil,
		buffer: buffer,
		lock:   sync.Mutex{},
	}

	return m
}

func (w *WSManager) RegisterWSConnection(ws net.Conn) {
	log.Println("new WS registration")
	newLive := &liveSocket{
		conn: ws,
	}
	w.lock.Lock()
	w.live = newLive
	w.lock.Unlock()

	go w.listenForDataFromWSUntilDead()
	go w.sendBuffer()
}

func (w *WSManager) listenForDataFromWSUntilDead() {
	for {
		msg, op, err := wsutil.ReadClientData(w.live.conn)
		if err != nil {
			log.Println("dead WS, exiting read loop ", msg, op, err)
			w.lock.Lock()
			w.live.conn.Close()
			w.live = nil
			w.lock.Unlock()
			break
		}
		log.Println("ws data: ", msg, op)
	}
}

func (w *WSManager) sendBuffer() {
	for w.live != nil {
		msg := <-w.buffer
		err := w.SendToSocket(msg)
		if err != nil {
			log.Println("dead WS, exiting send loop ", err, msg)
			w.lock.Lock()
			w.live.conn.Close()
			w.live = nil
			w.lock.Unlock()
			break
		}
	}
}

func (w *WSManager) SendToSocket(msg []byte) error {
	err := wsutil.WriteServerMessage(w.live.conn, ws.OpText, msg)

	if err != nil {
		log.Println("error writing to WS: ", err, msg)
	}

	return err
}
