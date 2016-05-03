package events

type Processor interface {
	Receive(*Event)
	Process(*Event)
	Send(*Event)
	Ack(*Event)
}
