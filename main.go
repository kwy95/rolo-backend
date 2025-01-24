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

	processedDataBuffer := make(chan []byte, 5000)
	arduinoSerial := arduino.NewArduinoSerial(processedDataBuffer)
	arduinoSerial.Start()

	accessAPI := api.NewAccessAPI(processedDataBuffer)
	accessAPI.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
