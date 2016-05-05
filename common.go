package events

// Event
type Key string
type Vals map[string]interface{}

type Event struct {
	Key  Key
	Vals Vals
}

func MakeEvent(k *Key, vals *Vals) *Event {
	return &Event{
		Key:  *k,
		Vals: *vals,
	}
}

// Bolt
type Bolt interface {
	Link(*Wire)
}

// Sender
type Sender interface {
	Bolt
	Send(*Event) error
}

type SenderBase struct {
	outlets []*Wire
}

func (s SenderBase) Link(wire *Wire) {
	s.outlets = append(s.outlets, wire)
}

func (s SenderBase) Send(evt *Event) error {
	for _, w := range s.outlets {
		*w.events <- evt
	}
	return nil
}

// Receiver
type Receiver interface {
	Bolt
	Receive(*Event) error
}

type ReceiverBase struct {
	inlets []*Wire
}

func (s ReceiverBase) Link(wire *Wire) {
	s.inlets = append(s.inlets, wire)
}
