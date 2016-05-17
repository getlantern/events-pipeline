package events

import (
	"fmt"
	"time"

	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("events-pipeline")
)

// Event
type Key string
type Vals map[string]interface{}

type Event struct {
	Key       Key
	Timestamp time.Time
	Vals      Vals

	// Internal
	wire   *Wire
	sender Sender
}

func NewEvent(k Key, vals *Vals) *Event {
	return &Event{
		Key:       k,
		Timestamp: time.Now(),
		Vals:      *vals,
	}
}

type SysEvent string

// System Events
// These do not have any key, and the ID is found as a Vals key
const (
	SystemEventInit SysEvent = "init"
	SystemEventStop SysEvent = "stop"
	SystemEventMark SysEvent = "mark"
)

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
	Send(*Event) error
}

type FeedbackFunc func(e *Event) error

type SenderBase struct {
	outlets         []*Wire
	feedbackHandler FeedbackFunc
}

func (s *SenderBase) ID() string {
	return ""
}

func (s *SenderBase) LinkOutlet(wire *Wire) {
	s.outlets = append(s.outlets, wire)
}

func (s *SenderBase) Send(evt *Event) error {
	for _, w := range s.outlets {
		copy := *evt
		copy.wire = w
		copy.sender = s
		*w.events <- &copy
	}
	return nil
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

func (s *ReceiverBase) Receive(evt *Event) error {
	return nil
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

func NewEmitterBase(ID string, feedback FeedbackFunc) *EmitterBase {
	return &EmitterBase{
		id:         ID,
		SenderBase: SenderBase{feedbackHandler: feedback},
	}
}

func (e *EmitterBase) ID() string {
	return e.id
}

func (e *EmitterBase) Emit(k Key, v *Vals) error {
	if k == "" {
		return fmt.Errorf("Event Key cannot be empty")
	}
	return e.Send(NewEvent(k, v))
}

// Sink
type Sink interface {
	Receiver
}

type SinkBase struct {
	ReceiverBase
	id string
}

func NewSinkBase(id string) *SinkBase {
	return &SinkBase{
		id: id,
	}
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

	// Receiver
	LinkInlet(*Wire)
	Receive(*Event) error
}

type ProcessorBase struct {
	id string
	SenderBase
	ReceiverBase
}

func NewProcessorBase(id string, feedback func(e *Event) error) *ProcessorBase {
	return &ProcessorBase{
		id:         id,
		SenderBase: SenderBase{feedbackHandler: feedback},
	}
}
