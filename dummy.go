package events

// Dummy Sink
type DummySink struct {
	*SinkBase
}

func NewDummySink(id string) *DummySink {
	return &DummySink{
		SinkBase: &SinkBase{
			id: id,
		},
	}
}

func (s *DummySink) Receive(evt *Event) error {
	log.Tracef("SINK ID %v received event: %v with: %v", s.ID(), evt.Key, evt.Vals)
	return evt.sender.Feedback(evt)
}

// Dummy Processor
type DummyProcessor struct {
	*ProcessorBase
}

func NewDummyProcessor(id string) *DummyProcessor {
	return &DummyProcessor{
		ProcessorBase: &ProcessorBase{
			id: id,
		},
	}
}

func (p *DummyProcessor) Receive(evt *Event) error {
	log.Tracef("PROCESSOR ID %v received event: %v with: %v", p.ID(), evt.Key, evt.Vals)

	err := p.Process(evt)
	if err != nil {
		return err
	}

	err = p.ProcessorBase.Feedback(evt)
	if err != nil {
		return err
	}

	_, err = p.Send(evt)

	return err
}

func (p *DummyProcessor) Send(evt *Event) (ack bool, err error) {
	for _, w := range p.outlets {
		copy := *evt
		copy.wire = w
		copy.sender = p
		*w.events <- &copy
	}
	return true, nil
}

func (p *DummyProcessor) Process(evt *Event) error {
	return nil
}

func (p *DummyProcessor) Feedback(evt *Event) error {
	log.Tracef("PROCESSOR ID %v received FEEDBACK of: %v with: %v", p.ID(), evt.Key, evt.Vals)
	return nil
}
