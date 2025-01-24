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

	arduino.Start()

	processedDataBuffer := make(chan []byte, 5000)
	accessAPI := api.NewAccessAPI(processedDataBuffer)
	accessAPI.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
