package events

type Bolt interface{}

type Sender interface {
	Bolt
	Send(e *Event) error
	Link(*Receiver)
}

type Receiver interface {
	Bolt
	Receive(e *Event) error
}
