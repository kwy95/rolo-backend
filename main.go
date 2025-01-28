package main

import (
	"log"
	"rolo/backend/api"
	"rolo/backend/arduino"
	"sync"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Println("working")

	directDataBuffer := make(chan []byte, 5000)

	arduinoSerial := arduino.NewArduinoSerial(directDataBuffer)
	arduinoSerial.Start()

	processedDataBuffer := make(chan []byte, 5000)

	accessAPI := api.NewAccessAPI(directDataBuffer, processedDataBuffer)
	accessAPI.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
