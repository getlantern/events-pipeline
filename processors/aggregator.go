package processors

import (
	events "github.com/getlantern/events-pipeline"
)

type AggregationDirective struct {
	Key            events.Key
	Val            string
	AggregatorFunc func(accum, x interface{}) (accum2, x2 interface{})
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
				a.currentValues[i], evt.Vals[d.Val] = d.AggregatorFunc(a.currentValues[i], val)
			}
		}
	}

	return a.ProcessorBase.Send(evt)
}

// Predefined identity values

var (
	RunningSumIdentity    = 0
	MovingAverageIdentity = []float64{0, 0}
)

// Predefined aggregator functions

func AggregatorIntRunningSum(accum, x interface{}) (accum2, x2 interface{}) {
	newSum := accum.(int) + x.(int)
	return newSum, newSum
}

func AggregatorFloat64RunningSum(accum, x interface{}) (accum2, x2 interface{}) {
	newSum := accum.(float64) + x.(float64)
	return newSum, newSum
}

func AggregatorFloat64MovingAverage(accum, x interface{}) (accum2, x2 interface{}) {
	currentMA := accum.([]float64)

	numElem := currentMA[0]
	newNumElem := numElem + 1

	currentAv := currentMA[1]

	newAverage := ((numElem * currentAv) + x.(float64)) / newNumElem
	newMovingAverage := []float64{newNumElem, newAverage}

	return newMovingAverage, newAverage
}
