package events

// Dummy Sink
type DummySink struct {
	ReceiverBase
}

func NewDummySink() *DummySink {
	return &DummySink{}
}

func (s *DummySink) Receive(evt *Event) error {
	log.Tracef("SINK Received event Key: %v Vals: %v", evt.Key, evt.Vals)
	return nil
}

// Dummy Processor
type DummyProcessor struct {
	ProcessorBase
}

func NewDummyProcessor() *DummyProcessor {
	return &DummyProcessor{}
}

func (p *DummyProcessor) Receive(evt *Event) error {
	log.Tracef("PROCESSOR Received event Key: %v Vals: %v", evt.Key, evt.Vals)
	return p.Send(evt)
}
