package racer

import (
	"log"
	"strconv"
)

type Race struct {
	inputBuff  chan []byte
	outputBuff chan []byte
}

type Bike struct {
	ID       int `json:"id"`
	Speed    int `json:"speed"`
	Distance int `json:"distance"`
}

type BikeInstant struct {
	ID          int `json:"id"`
	Moment      int `json:"moment"`
	PulseLength int `json:"pulseLength"`
}

func NewRace(in chan []byte, out chan []byte) *Race {
	return &Race{
		inputBuff:  in,
		outputBuff: out,
	}
}

func (r *Race) Start() {
	go func() {
		for {
			r.outputBuff <- <-r.inputBuff
		}
	}()
}

func bikeDataFromArduinoData(data []byte) *BikeInstant {
	whichField := 0
	counter := 0
	buffer := make([]byte, 16)
	bikeInstant := BikeInstant{}

	for i := 1; i < len(data); i++ {
		if data[i] == ';' || data[i] == '}' {
			err := bikeInstant.inputField(whichField, buffer[:counter])
			if err != nil {
				log.Println("error parsing bike data: ", err)
				return nil
			}
			whichField++
			counter = 0
			continue
		}
		buffer[counter] = data[i]
		counter++
	}

	return &bikeInstant
}

func (b *BikeInstant) inputField(field int, value []byte) error {
	var err error
	switch field {
	case 0:
		b.ID, err = strconv.Atoi(string(value))
	case 1:
		b.Moment, err = strconv.Atoi(string(value))
	case 2:
		b.PulseLength, err = strconv.Atoi(string(value))
	}
	return err
}
