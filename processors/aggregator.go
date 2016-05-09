package processors

import (
	events "github.com/getlantern/events-pipeline"
	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("events-pipeline")
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
	log.Tracef("AGGREGATOR ID %v received event: %v with: %v", a.ID(), evt.Key, evt.Vals)

	err := a.Process(evt)
	if err != nil {
		return err
	}

	// Should always be ok
	/*
		err = a.ProcessorBase.Feedback(evt)
		if err != nil {
			return err
		}
	*/

	return err
}

func (a *Aggregator) Send(evt *events.Event) (ack bool, err error) {
	return a.ProcessorBase.Send(evt)
}

func (a *Aggregator) Process(evt *events.Event) error {
	// TODO: here is where the actual aggregation happens
	_, err := a.Send(evt)

	return err
}
