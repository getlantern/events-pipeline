package events

type Processor interface {
	Send(*Event) error
	Receive(*Event) error
	Process(*Event)
	Ack(*Event)
}

type DummyProcessor struct {
	SenderBase
	ReceiverBase
}

func MakeDummyProcessor() *DummyProcessor {
	return &DummyProcessor{}
}

func (e *DummyProcessor) Send(k Key, v *Vals) error {
	return e.SenderBase.Send(MakeEvent(k, v))
}
