package events

type Sink interface {
	Receive(*Event)
}
