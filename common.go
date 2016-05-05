package events

import (
	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("events-pipeline")
)

// Event
type Key string
type Vals map[string]interface{}

type Event struct {
	Key  Key
	Vals Vals
}

func MakeEvent(k Key, vals *Vals) *Event {
	return &Event{
		Key:  k,
		Vals: *vals,
	}
}

// Bolt
type Bolt interface {
	Link(*Wire)
}

// Wire
type Wire struct {
	senders   []Sender
	receivers []Receiver
	events    *chan *Event
}

// Sender
type Sender interface {
	Link(*Wire)
	Send(*Event) error
}

type SenderBase struct {
	outlets []*Wire
}

func (s *SenderBase) Link(wire *Wire) {
	s.outlets = append(s.outlets, wire)
}

func (s *SenderBase) Send(evt *Event) error {
	for _, w := range s.outlets {
		*w.events <- evt
	}
	return nil
}

// Receiver
type Receiver interface {
	Link(*Wire)
	Receive(*Event) error
}

type ReceiverBase struct {
	inlets []*Wire
}

func (s *ReceiverBase) Link(wire *Wire) {
	s.inlets = append(s.inlets, wire)
}
