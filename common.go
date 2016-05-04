package events

type Receiver interface {
	Receive(*Event) error
}

type Sender interface {
	Send() error
	Link(*Receiver)
}
