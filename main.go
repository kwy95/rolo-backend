package main

import (
	"log"
	"rolo/backend/arduino"
)

func main() {
	log.Println("working")
	arduino.Start()
}
