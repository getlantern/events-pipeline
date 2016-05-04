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
type Bolt interface{}

// Sender
type Sender interface {
	Send(*Event) error
}

type SenderBase struct {
	outlets []*Wire
}

func (s SenderBase) Send(evt *Event) error {
	// TODO
	return nil
}

// Receiver
type Receiver interface {
	Receive(*Event) error
}

type ReceiverBase struct {
	inlets []*Wire
}

func (s ReceiverBase) Receive(evt *Event) error {
	// TODO
	return nil
}
