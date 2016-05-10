package processors

import (
	events "github.com/getlantern/events-pipeline"
)

type AggregatorOptions struct {
	FeedbackHandler func(e *events.Event) error
}

type Aggregator struct {
	*events.ProcessorBase

	options *AggregatorOptions
}

func NewAggregator(id string, feedback events.FeedbackFunc, opts *AggregatorOptions) *Aggregator {
	return &Aggregator{
		ProcessorBase: events.NewProcessorBase(id, feedback),
		options:       opts,
	}
}

func (a *Aggregator) Receive(evt *events.Event) error {
	log.Tracef("AGGREGATOR ID %v PROCESSED event: %v with: %v", a.ID(), evt.Key, evt.Vals)
	// TODO Processing

	return a.ProcessorBase.Receive(evt)
}

func (a *Aggregator) Send(evt *events.Event) error {
	return a.ProcessorBase.Send(evt)
}
