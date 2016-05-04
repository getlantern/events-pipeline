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
	wires []*Wire
}

func NewPipeline(sender Sender) *Pipeline {
	return &Pipeline{
		bolts: []Bolt{sender},
		wires: []*Wire{},
	}
}

func (p *Pipeline) Plug(s Sender, r Receiver) error {
	newWire := &Wire{
		sender:   &s,
		receiver: &r,
	}
	p.wires = append(p.wires, newWire)
	s.Link(newWire)
	r.Link(newWire)

	return nil
}

func (p *Pipeline) Run() {
	for _, wire := range p.wires {
		receiver := *wire.receiver
		go func() {
			for {
				select {
				case evt := <-*wire.events:
					receiver.Receive(evt)
					fmt.Println(evt)
				}
			}
		}()
	}
}
