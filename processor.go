package events

type Processor interface {
	Send(*Event) error
	Receive(*Event) error
	Process(*Event)
	Ack(*Event)
}
