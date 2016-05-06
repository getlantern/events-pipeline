package events

// Dummy Sink
type DummySink struct {
	ReceiverBase
	id string
}

func NewDummySink(id string) *DummySink {
	return &DummySink{id: id}
}

func (s *DummySink) Receive(evt *Event) error {
	log.Tracef("SINK Received event Key: %v Vals: %v", evt.Key, evt.Vals)
	return nil
}

func (s *DummySink) ID() string {
	return s.id
}

// Dummy Processor
type DummyProcessor struct {
	ProcessorBase
	id string
}

func NewDummyProcessor(id string) *DummyProcessor {
	return &DummyProcessor{id: id}
}

func (p *DummyProcessor) Receive(evt *Event) error {
	log.Tracef("PROCESSOR Received event Key: %v Vals: %v", evt.Key, evt.Vals)
	return p.Send(evt)
}

func (p *DummyProcessor) ID() string {
	return p.id
}
