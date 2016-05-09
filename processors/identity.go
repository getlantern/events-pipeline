package processors

import (
	events "github.com/getlantern/events-pipeline"
)

// Identity Processor
type IdentityProcessor struct {
	*events.ProcessorBase
}

func NewIdentityProcessor(id string) *IdentityProcessor {
	return &IdentityProcessor{
		ProcessorBase: events.NewProcessorBase(id),
	}
}

func (p *IdentityProcessor) Receive(evt *events.Event) error {
	log.Tracef("PROCESSOR ID %v received event: %v with: %v", p.ID(), evt.Key, evt.Vals)
	return p.ProcessorBase.Receive(evt)
}

func (p *IdentityProcessor) Process(evt *events.Event) error {
	// Do nothing
	return nil
}

func (p *IdentityProcessor) Feedback(evt *events.Event) error {
	log.Tracef("PROCESSOR ID %v received FEEDBACK of: %v with: %v", p.ID(), evt.Key, evt.Vals)
	return p.ProcessorBase.Feedback(evt)
}
