package processors

import (
	events "github.com/getlantern/events-pipeline"
)

type AggregatorOptions struct {
}

type Aggregator struct {
	*events.ProcessorBase

	options *AggregatorOptions
}

func NewAggregator(id string, opts *AggregatorOptions) *Aggregator {
	return &Aggregator{
		ProcessorBase: events.NewProcessorBase(id),
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
