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
	done       bool
	Bikes      map[int]*Bike `json:"bikes"`
	Order      int           `json:"order"`
}

type Bike struct {
	LastUpdate   uint32  `json:"-"`
	SpeedOn      float64 `json:"speed"`
	SpeedBetween float64 `json:"-"`
	Distance     float64 `json:"-"`
	Progress     float64 `json:"progress"`
}

type BikeInstant struct {
	ID          int    `json:"id"`
	Moment      uint32 `json:"moment"`
	PulseLength uint32 `json:"pulseLength"`
}

func NewRace(in chan []byte, out chan []byte) *Race {
	race := Race{
		inputBuff:  in,
		outputBuff: out,
		done:       false,
		Bikes:      make(map[int]*Bike),
	}
	bike0 := Bike{}
	bike1 := Bike{}
	race.Bikes[0] = &bike0
	race.Bikes[1] = &bike1

	return &race
}

func (r *Race) Start() {
	clearChannel(r.inputBuff)
	go func() {
		for data := range r.inputBuff {
			if r.Bikes[0].Progress >= 1.0 || r.Bikes[1].Progress >= 1.0 || r.done == true {
				log.Println("Finished:", r.Bikes[0].Progress, ";", r.Bikes[1].Progress)
				return
			}

			update := bikeDataFromArduinoData(data)
			r.processUpdate(update)
			encoded, _ := json.Marshal(r)
			r.outputBuff <- encoded
		}
	}()
}

func (r *Race) EndRace() {
	log.Println("Race finished:", r.Bikes[0].Progress, ";", r.Bikes[1].Progress)
	r.done = true
}

func clearChannel(c chan []byte) {
	for len(c) > 0 {
		data := <-c
		func(d []byte) {}(data)
	}
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

func (r *Race) processUpdate(update *BikeInstant) {
	b, ok := r.Bikes[update.ID]
	if !ok {
		log.Println("bike not found")
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
	b.Progress = b.Distance / TrackLength
	r.Order++
}

func (b *BikeInstant) inputField(field int, value []byte) error {
	var err error
	var convertedValue uint64
	switch field {
	case 0:
		b.ID, err = strconv.Atoi(string(value))
	case 1:
		convertedValue, err = strconv.ParseUint(string(value), 10, 32)
		b.PulseLength = uint32(convertedValue)
	case 2:
		convertedValue, err = strconv.ParseUint(string(value), 10, 32)
		b.Moment = uint32(convertedValue)
	}
	return err
}
