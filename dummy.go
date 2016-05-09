package events

// Null Sink
type NullSink struct {
	*SinkBase
}

func NewNullSink(id string) *NullSink {
	return &NullSink{
		SinkBase: &SinkBase{
			id: id,
		},
	}
}

func (s *NullSink) Receive(evt *Event) error {
	log.Tracef("SINK ID %v received event: %v with: %v", s.ID(), evt.Key, evt.Vals)
	return s.SinkBase.Receive(evt)
}

// Identity Processor
type IdentityProcessor struct {
	*ProcessorBase
}

func NewIdentityProcessor(id string) *IdentityProcessor {
	return &IdentityProcessor{
		ProcessorBase: &ProcessorBase{
			id: id,
		},
	}
}

func (p *IdentityProcessor) Receive(evt *Event) error {
	log.Tracef("PROCESSOR ID %v received event: %v with: %v", p.ID(), evt.Key, evt.Vals)
	return p.ProcessorBase.Receive(evt)
}

func (p *IdentityProcessor) Process(evt *Event) error {
	// Do nothing
	return nil
}

func (p *IdentityProcessor) Feedback(evt *Event) error {
	log.Tracef("PROCESSOR ID %v received FEEDBACK of: %v with: %v", p.ID(), evt.Key, evt.Vals)
	return p.ProcessorBase.Feedback(evt)
}
