package events

import (
	"fmt"
)

type Pipeline struct {
	Bolts []Bolt
	Wires []*Wire
	stop  chan struct{}
}

func NewPipeline(sender Sender) *Pipeline {
	return &Pipeline{
		Bolts: []Bolt{sender},
		Wires: []*Wire{},
		stop:  make(chan struct{}),
	}
}

func (p *Pipeline) Plug(s Sender, r Receiver) (*Wire, error) {
	evChan := make(chan *Event)
	newWire := &Wire{
		senders:   []Sender{s},
		receivers: []Receiver{r},
		events:    &evChan,
	}
	p.Wires = append(p.Wires, newWire)
	s.Link(newWire)
	r.Link(newWire)

	return newWire, nil
}

func (p *Pipeline) PlugWith(s Sender, r Receiver, wire *Wire) error {
	// Find out if the sender already uses this wire
	var fSender Sender
	for _, ws := range wire.senders {
		if ws == s {
			fSender = ws
			break
		}
	}
	var fReceiver Receiver
	for _, wr := range wire.receivers {
		if wr == r {
			fReceiver = wr
			break
		}
	}
	if fSender != nil && fReceiver != nil {
		return fmt.Errorf("This wire already connects these bolts")
	}

	return nil
}

func (p *Pipeline) Run() {
	for _, wire := range p.Wires {
		go func() {
			for {
				select {
				case evt := <-*wire.events:
					for _, rcv := range wire.receivers {
						rcv.Receive(evt)
					}

					log.Tracef("Received event; %v", evt)
				case <-p.stop:
					return
				}
			}
		}()
	}
}

func (p *Pipeline) Stop() {
	p.stop <- struct{}{}
}
