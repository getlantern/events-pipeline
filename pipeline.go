package events

import (
	"fmt"

	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("events-pipeline")
)

type EventChannel chan (*Event)

type Wire struct {
	senders   []Sender
	receivers []Receiver
	events    *EventChannel
}

type Pipeline struct {
	Bolts []Bolt
	Wires []*Wire
}

func NewPipeline(sender Sender) *Pipeline {
	return &Pipeline{
		Bolts: []Bolt{sender},
		Wires: []*Wire{},
	}
}

func (p *Pipeline) Plug(s Sender, r Receiver) (*Wire, error) {
	newWire := &Wire{
		senders:   []Sender{s},
		receivers: []Receiver{r},
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
				}
			}
		}()
	}
}
