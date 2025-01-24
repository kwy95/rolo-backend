package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

func main() {
	echoMessages()
}

func echoMessages() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{})
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})

	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("error reading settlements msg:", err)
			break
		}
		// je := json.Unmarshal(message, &raceUpdate)
		// if je != nil {
		// 	log.Println("parse err", je)
		// 	break
		// }
		log.Println("msg:", string(message))
	}
	log.Println("falling out of subs method")
}
