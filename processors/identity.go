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
	err := p.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	// Processing could be done here

	return p.ProcessorBase.Send(evt)
}

func (p *IdentityProcessor) Feedback(evt *events.Event) error {
	return p.ProcessorBase.Feedback(evt)
}
