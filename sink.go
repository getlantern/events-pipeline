package events

import (
	"fmt"
)

type Sink interface {
	Receiver
}

type DummySink struct {
	ReceiverBase
}

func NewDummySink() *DummySink {
	return &DummySink{}
}

func (s *DummySink) Link(w *Wire) {
	s.ReceiverBase.Link(w)
}

func (s *DummySink) Receive(evt *Event) error {
	fmt.Printf("Received event Key: %v Vals: %v\n", evt.Key, evt.Vals)
	return nil
}
