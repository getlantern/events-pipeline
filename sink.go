package events

import (
	"fmt"
)

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
