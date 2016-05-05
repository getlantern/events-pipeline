package events

import (
	"fmt"
)

// Dummy Emitter

// Dummy Sink

type DummySink struct {
	ReceiverBase
}

func NewDummySink() *DummySink {
	return &DummySink{}
}

/*
func (s *DummySink) LinkInlet(w *Wire) {
	s.ReceiverBase.LinkInlet(w)
}
*/
func (s *DummySink) Receive(evt *Event) error {
	fmt.Printf("Received event Key: %v Vals: %v\n", evt.Key, evt.Vals)
	return nil
}

// Dummy Processor

type DummyProcessor struct {
	ProcessorBase
}

func NewDummyProcessor() *DummyProcessor {
	return &DummyProcessor{}
}

func (p *DummyProcessor) LinkOut(w *Wire) {
	p.SenderBase.LinkOutlet(w)
}

func (p *DummyProcessor) LinkInlet(w *Wire) {
	p.ReceiverBase.LinkInlet(w)
}

func (p *DummyProcessor) Receive(evt *Event) error {
	return p.Send(evt)
}

func (p *DummyProcessor) Send(evt *Event) error {
	return p.ProcessorBase.Send(evt)
}
