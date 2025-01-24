package racer

type Race struct {
	inputBuff  chan []byte
	outputBuff chan []byte
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
