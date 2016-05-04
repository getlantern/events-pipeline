package events

import (
	"fmt"
)

type EventChannel chan (*Event)

type Wire struct {
	sender   *Sender
	receiver *Receiver
	events   *EventChannel
}

type Pipeline struct {
	bolts []Bolt
	wires []Wire
}

func NewPipeline(emitter *Emitter) *Pipeline {
	return &Pipeline{
		bolts: []Bolt{emitter},
		wires: []Wire{},
	}
}

func (p *Pipeline) Plug(s *Sender, r *Receiver) error {
	newWire := Wire{
		sender:   s,
		receiver: r,
	}
	p.wires = append(p.wires, newWire)

	return nil
}

func (p *Pipeline) Run() {
	for _, wire := range p.wires {
		go func() {
			for {
				select {
				case evt := <-*wire.events:
					// TODO: Dispatch depending on type
					fmt.Println(evt)
				}
			}
		}()
	}
}
