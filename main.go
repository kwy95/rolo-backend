package main

import (
	"log"
	"rolo/backend/arduino"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Println("working")
	arduino.Start()
}
