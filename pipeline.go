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
	Bolts map[string]Bolt
	Wires []*Wire

	init chan struct{}
	stop chan struct{}
}

func NewPipeline(sender Sender) *Pipeline {
	p := &Pipeline{
		Bolts: map[string]Bolt{sender.ID(): sender},
		Wires: []*Wire{},
		init:  make(chan struct{}),
		stop:  make(chan struct{}),
	}
	go func() {
		p.init <- struct{}{}
	}()
	return p
}

func (p *Pipeline) Plug(s Sender, r Receiver) (*Wire, error) {
	evChan := make(chan *Event)
	wire := &Wire{
		senders:   []Sender{},
		receivers: []Receiver{},
		events:    &evChan,
	}

	w, err := p.PlugWith(s, r, wire)
	if err == nil {
		p.Wires = append(p.Wires, wire)
	}
	return w, err
}

func (p *Pipeline) PlugWith(s Sender, r Receiver, wire *Wire) (*Wire, error) {
	if _, exists := p.Bolts[s.ID()]; !exists {
		p.Bolts[s.ID()] = s
	}
	if _, exists := p.Bolts[r.ID()]; !exists {
		p.Bolts[r.ID()] = r
	}

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
	// Initialization
	if err := p.broadcastSysEvent(SystemEventInit); err != nil {
		log.Errorf("Error broadcasting INIT system event: %v", err) //
	}

	for _, wire := range p.Wires {
		// Copy the reference to the wire
		// Remove this line and you will unleash the wrath of the gods
		wire := wire
		go func() {
			for {
				// We always select on the stop signal at the end, so we make sure
				// all wires are cleared of events before stopping.
				// This only works if the system has bounded liveness and cannot lock
				// indefinitely
				select {
				case evt := <-*wire.events:
					// Processing
					for _, rcv := range wire.receivers {
						err := rcv.Receive(evt)
						if err != nil {
							log.Errorf("Error receiving event: %v", err)
							continue
						}
						if evt.sender.(*SenderBase).feedbackHandler != nil {
							err = evt.sender.(*SenderBase).feedbackHandler(evt)
							if err != nil {
								log.Errorf("Error in feedback handler: %v", err)
							}

						}
					}
				case <-p.stop:
					return
				}
			}
		}()
	}
}

func (p *Pipeline) Stop() {
	// Broadcast a system event to all receivers
	if err := p.broadcastSysEvent(SystemEventStop); err != nil {
		log.Errorf("Error broadcasting STOP system event: %v", err)
	}

	p.stop <- struct{}{}
}

func (p *Pipeline) broadcastSysEvent(sysEvType SysEvent) error {
	for _, b := range p.Bolts {
		if r, ok := b.(Receiver); ok {
			err := r.Receive(NewEvent("", &Vals{string(sysEvType): nil}))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
