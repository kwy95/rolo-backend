package racer

import (
	"encoding/json"
	"log"
	"strconv"
)

const (
	// meters
	TrackLength  float64 = 200.0
	WheelRadius  float64 = 0.349
	SensorRadius float64 = 0.12
	SensorArch   float64 = 0.03

	BikeRatio   float64 = 44.0 / 14.0
	PulsePerRev float64 = 1.0
	TimeScale   float64 = 1e-6 // 1 microsecond
	MpStoKmH    float64 = 3.6
	Pi          float64 = 3.14159265358979323846264338327950288419716939937510582097494459

	DistBetweenPulse float64 = 2.0 * Pi * WheelRadius * BikeRatio / PulsePerRev
	DistOnPulse      float64 = WheelRadius * SensorArch * BikeRatio / SensorRadius
)

type Race struct {
	inputBuff  chan []byte
	outputBuff chan []byte
	Bikes      [2]*Bike
}

type Bike struct {
	ID           int     `json:"id"`
	Order        int     `json:"order"`
	LastUpdate   int     `json:"-"`
	SpeedOn      float64 `json:"speed"`
	SpeedBetween float64 `json:"-"`
	Distance     float64 `json:"distance"`
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
		Bikes:      [2]*Bike{{ID: 0}, {ID: 1}},
	}
}

func (r *Race) Start() {
	go func() {
		for data := range r.inputBuff {
			update := bikeDataFromArduinoData(data)
			updatedBike := r.Bikes[update.ID]
			updatedBike.processUpdate(update)
			encoded, _ := json.Marshal(updatedBike)
			r.outputBuff <- encoded
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

func (b *Bike) processUpdate(update *BikeInstant) {
	if b.ID != update.ID {
		log.Println("bike ID mismatch")
		return
	}

	if b.LastUpdate == 0 {
		b.LastUpdate = update.Moment
		return
	}

	timeDiff := update.Moment - b.LastUpdate
	b.LastUpdate = update.Moment

	b.SpeedOn = MpStoKmH * DistOnPulse / (float64(update.PulseLength) * TimeScale)
	b.SpeedBetween = MpStoKmH * DistBetweenPulse / (float64(timeDiff) * TimeScale)
	b.Distance += DistBetweenPulse
	b.Order++
}

func (b *BikeInstant) inputField(field int, value []byte) error {
	var err error
	switch field {
	case 0:
		b.ID, err = strconv.Atoi(string(value))
	case 1:
		b.PulseLength, err = strconv.Atoi(string(value))
	case 2:
		b.Moment, err = strconv.Atoi(string(value))
	}
	return err
}
