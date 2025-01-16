package arduino

import (
	"log"
	"time"
	"go.bug.st/serial"
)

func Start() {
	log.Println("started")

	ports := getPorts()
	port := connectArduino(ports)
	if port == nil {
		log.Fatal("Failed to connect to Arduino")
	}
	defer port.Close()
}

func getPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	return ports
}

func connectArduino(ports []string) serial.Port {
	mode := &serial.Mode{
		BaudRate: 9600,
	}

	for _, port_name := range ports {
		port, err := serial.Open(port_name, mode)
		if err != nil {
			log.Fatal(err)
		}

		if !confirmArduino(port) {
			port.Close()
			continue
		}

		return port
	}

	return nil
}

func confirmArduino(port serial.Port) bool {
	port.ResetInputBuffer()
	time.Sleep(1)
	buff := make([]byte, 1)
	n, err := port.Read(buff)
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		return false
	}

	if string(buff[:n]) == "A" {
		n, err = port.Write([]byte("A"))
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			log.Fatal("Failed to write to Arduino")
		}
		log.Printf("sent %v bytes to arduino\n", n)
		port.ResetInputBuffer()
		return true
	}

	return false
}
