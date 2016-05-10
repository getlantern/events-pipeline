package processors

import (
	events "github.com/getlantern/events-pipeline"
)

type AggregationDirective struct {
	Key            events.Key
	Val            string
	AggregatorFunc func(accum, x interface{}) interface{}
	Identity       interface{}
}

type Aggregator struct {
	*events.ProcessorBase

	directives    []AggregationDirective
	currentValues []interface{}
}

func NewAggregator(id string, feedback events.FeedbackFunc, ds ...AggregationDirective) *Aggregator {
	initValues := make([]interface{}, len(ds))
	for i, d := range ds {
		initValues[i] = d.Identity
	}

	return &Aggregator{
		ProcessorBase: events.NewProcessorBase(id, feedback),
		currentValues: initValues,
		directives:    ds,
	}
}

func (a *Aggregator) Receive(evt *events.Event) error {
	log.Tracef("AGGREGATOR ID %v PROCESSED event: %v with: %v", a.ID(), evt.Key, evt.Vals)
	err := a.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	for i, d := range a.directives {
		if evt.Key == d.Key {
			if val, ok := evt.Vals[d.Val]; ok {
				updated := d.AggregatorFunc(a.currentValues[i], val)
				a.currentValues[i] = updated
				evt.Vals[d.Val] = updated
			}
		}
	}

	return a.ProcessorBase.Send(evt)
}

// Predefined aggregator functions

func AggregatorIntRunningSum(accum, x interface{}) interface{} {
	return accum.(int) + x.(int)
}
