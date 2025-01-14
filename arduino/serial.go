package arduino

import (
	"log"
	"go.bug.st/serial"
)

func Start() {
	log.Println("started")
	getPorts()
}

func getPorts() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		log.Printf("Found port: %v\n", port)
	}
}
