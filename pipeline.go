// The pipeline is the engine of the events processing mechanism.
// It is implemented as a set of of "wires" that concurrently de/mux
// the data from "bolts" (senders or receivers of data).
// Data processing is synchronized in the receivers, which means that
// they need to be ready to read from the wire at all times.
// This is particularly relevant for sinks, which in all likelyhood
// will rely on networking (sending processed events to a remote
// service). In practical terms, this means that they have to implement
// their own buffering scheme internally so the pipeline implementation
// is more lightweight and we have more control of the buffering
// depending on the sink requirements.

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
		Bolts: []Bolt{sender}, // TODO
		Wires: []*Wire{},
		stop:  make(chan struct{}),
	}
}

func (p *Pipeline) Plug(s Sender, r Receiver) (*Wire, error) {
	evChan := make(chan *Event)
	wire := &Wire{
		senders:   []Sender{},
		receivers: []Receiver{},
		events:    &evChan,
	}
	p.Wires = append(p.Wires, wire)

	return p.PlugWith(s, r, wire)
}

func (p *Pipeline) PlugWith(s Sender, r Receiver, wire *Wire) (*Wire, error) {
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
		return nil, fmt.Errorf("This wire already connects these bolts")
	}

	wire.senders = append(wire.senders, s)
	wire.receivers = append(wire.receivers, r)
	s.LinkOutlet(wire)
	r.LinkInlet(wire)

	return wire, nil
}

func (p *Pipeline) Run() {
	for _, wire := range p.Wires {
		// Copy the reference to the wire
		// Remove this line and you will unleash the wrath of the gods
		wire := wire
		go func() {
			for {
				select {
				case evt := <-*wire.events:
					for _, rcv := range wire.receivers {
						rcv.Receive(evt)
					}
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
