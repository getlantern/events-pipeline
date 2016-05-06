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
	log.Tracef("SINK ID %v received event: %v with: %v", s.ID(), evt.Key, evt.Vals)
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
	log.Tracef("PROCESSOR ID %v received event: %v with: %v", p.ID(), evt.Key, evt.Vals)
	return p.Send(evt)
}

func (p *DummyProcessor) ID() string {
	return p.id
}
