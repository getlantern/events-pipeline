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

	// Internal
	wire   *Wire
	sender Sender
}

func NewEvent(k Key, vals *Vals) *Event {
	return &Event{
		Key:  k,
		Vals: *vals,
	}
}

// Bolt
type Bolt interface {
	ID() string
}

// Wire
type Wire struct {
	senders   []Sender
	receivers []Receiver
	events    *chan *Event
}

// Sender
type Sender interface {
	Bolt
	LinkOutlet(*Wire)
	Send(*Event) (ack bool, err error)
	Feedback(*Event) error
}

type SenderBase struct {
	outlets []*Wire
}

func (s *SenderBase) ID() string {
	return ""
}

func (s *SenderBase) LinkOutlet(wire *Wire) {
	s.outlets = append(s.outlets, wire)
}

// Receiver
type Receiver interface {
	Bolt
	LinkInlet(*Wire)
	Receive(*Event) error
}

type ReceiverBase struct {
	inlets []*Wire
}

func (s *ReceiverBase) LinkInlet(wire *Wire) {
	s.inlets = append(s.inlets, wire)
}

// Emitter
type Emitter interface {
	Sender
	Emit(Key, *Vals) error
}

type EmitterBase struct {
	id string
	SenderBase
}

func NewEmitter(ID string) *EmitterBase {
	return &EmitterBase{
		id: ID,
	}
}

func (e *EmitterBase) ID() string {
	return e.id
}

func (e *EmitterBase) Emit(k Key, v *Vals) error {
	_, err := e.Send(NewEvent(k, v))
	return err
}

func (e *EmitterBase) Send(evt *Event) (ack bool, err error) {
	for _, w := range e.outlets {
		copy := *evt
		copy.wire = w
		copy.sender = e
		*w.events <- &copy
	}
	return true, nil
}

func (s *SenderBase) Feedback(evt *Event) error {
	return nil
}

// Sink
type Sink interface {
	Receiver
}

type SinkBase struct {
	ReceiverBase
	id string
}

func (s *SinkBase) ID() string {
	return s.id
}

// Processor
type Processor interface {
	Bolt

	// Sender
	LinkOutlet(*Wire)
	Send(*Event) error
	Feedback(*Event) error

	// Receiver
	LinkInlet(*Wire)
	Receive(*Event) error

	Process(*Event) error
}

type ProcessorBase struct {
	id string
	SenderBase
	ReceiverBase
}

func (p *ProcessorBase) ID() string {
	return p.id
}
