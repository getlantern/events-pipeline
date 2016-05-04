package events

type Processor interface {
	Receiver
	Sender
	Process(*Event)
	Ack(*Event)
}
