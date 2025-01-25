package arduino

import (
	"go.bug.st/serial"
	"log"
	"time"
)

type ArduinoSerial struct {
	port   serial.Port
	buffer chan []byte
}

func NewArduinoSerial(buffer chan []byte) *ArduinoSerial {
	return &ArduinoSerial{
		port:   nil,
		buffer: buffer,
	}
}

func (a *ArduinoSerial) Start() {
	log.Println("started arduino")

	ports := getPorts()
	a.connectArduino(ports)

	go a.receiveData()
}

func (a *ArduinoSerial) receiveData() {
	defer func() {
		a.port.Close()
		a.port = nil
	}()

	buff := make([]byte, 100)
	index := 0
	for {
		n, err := a.port.Read(buff[index:])
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			log.Println("\nEOF")
			break
		}
		index = a.parseData(buff, index, n)
	}
}

func (a *ArduinoSerial) parseData(data []byte, index int, n int) int {
	for i := 0; i < index+n; i++ {
		if data[i] == '}' {
			a.buffer <- data[:i+1]
			for j := i + 1; j < index+n; j++ {
				data[j-i-1] = data[j]
			}
			return index + n - i - 1
		}
	}
	return index + n
}

func getPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	for len(ports) == 0 {
		ports, err = serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}
	}
	return ports
}

func (a *ArduinoSerial) connectArduino(ports []string) {
	mode := &serial.Mode{
		BaudRate: 115200,
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

		a.port = port
		return
	}

	log.Fatal("Failed to connect to Arduino")
}

func confirmArduino(port serial.Port) bool {
	port.ResetInputBuffer()
	time.Sleep(1 * time.Second)
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
