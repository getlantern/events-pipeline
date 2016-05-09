package sinks

import (
	events "github.com/getlantern/events-pipeline"
)

// Null Sink
type NullSink struct {
	*events.SinkBase
}

func NewNullSink(id string) *NullSink {
	return &NullSink{
		SinkBase: events.NewSinkBase(id),
	}
}

func (s *NullSink) Receive(evt *events.Event) error {
	return s.SinkBase.Receive(evt)
}
