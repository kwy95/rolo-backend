package main

import (
	"log"
	"rolo/backend/arduino"
	"sync"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Println("working")

	arduino.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
