package events

import (
	"testing"
	"time"

	"github.com/getlantern/testify/assert"
)

// Identity Processor
type IdentityProcessor struct {
	*ProcessorBase
}

func NewIdentityProcessor(id string, feedback func(e *Event) error) *IdentityProcessor {
	return &IdentityProcessor{
		ProcessorBase: NewProcessorBase(id, feedback),
	}
}

func (p *IdentityProcessor) Receive(evt *Event) error {
	if evt.Key == "" {
		return nil
	}

	log.Tracef("PROCESSOR ID %v received event: %v with: %v", p.ID(), evt.Key, evt.Vals)
	err := p.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	return p.ProcessorBase.Send(evt)
}

// Null Sink
type NullSink struct {
	*SinkBase
}

func NewNullSink(id string) *NullSink {
	return &NullSink{
		SinkBase: NewSinkBase(id),
	}
}

func (s *NullSink) Receive(evt *Event) error {
	log.Tracef("SINK ID %v received event: %v with: %v", s.ID(), evt.Key, evt.Vals)
	return s.SinkBase.Receive(evt)
}

func TestTrivialPipeline(t *testing.T) {
	t.SkipNow()

	emitter := NewEmitterBase("test-emitter", nil)
	sink := NewNullSink("test-sink")
	pipeline := NewPipeline(emitter)
	_, err := pipeline.Plug(emitter, sink)
	assert.Nil(t, err, "Should be nil")

	assert.Equal(t, 1, len(emitter.outlets))
	assert.Equal(t, 1, len(sink.inlets))

	pipeline.Run()

	emitter.Emit("Key A", &Vals{})
	emitter.Emit("Key B", &Vals{})
	time.Sleep(time.Millisecond * 20)

	pipeline.Stop()
}

func TestFeedback(t *testing.T) {
	quote1 := "Everything you can imagine is real"
	quote2 := "Life is wasted on the living"

	efc := make(chan string, 2)
	ipfc := make(chan string, 2)

	emitter := NewEmitterBase("test-emitter", func(e *Event) error {
		log.Tracef("FEEDBACK!")
		efc <- quote1
		return nil
	})
	sink := NewNullSink("test-sink")
	dummy := NewIdentityProcessor("test-processor", func(e *Event) error {
		ipfc <- quote2
		return nil
	})
	pipeline := NewPipeline(emitter)

	_, err := pipeline.Plug(emitter, dummy)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(dummy, sink)
	assert.Nil(t, err, "Should be nil")

	pipeline.Run()

	emitter.Emit("Key A", &Vals{})
	emitter.Emit("Key B", &Vals{})

	time.Sleep(time.Millisecond * 20)

	assert.Equal(t, 2, len(efc), "The channel should be filled with two elements")
	assert.Equal(t, 2, len(ipfc), "The channel should be filled with two elements")

	assert.Equal(t, quote1, <-efc, "The channel should contain the proper data")
	assert.Equal(t, quote1, <-efc, "The channel should contain the proper data")
	assert.Equal(t, quote2, <-ipfc, "The channel should contain the proper data")
	assert.Equal(t, quote2, <-ipfc, "The channel should contain the proper data")
	assert.Empty(t, efc, "The channel should be empty")
	assert.Empty(t, ipfc, "The channel should be empty")

	pipeline.Stop()
}

func TestProcessorChain(t *testing.T) {
	emitter := NewEmitterBase("test-emitter", nil)
	sink := NewNullSink("test-sink")
	dummy1 := NewIdentityProcessor("test-processor A", nil)
	dummy2 := NewIdentityProcessor("test-processor B", nil)
	dummy3 := NewIdentityProcessor("test-processor C", nil)
	pipeline := NewPipeline(emitter)
	_, err := pipeline.Plug(emitter, dummy1)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(dummy1, dummy2)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(dummy2, dummy3)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(dummy3, sink)
	assert.Nil(t, err, "Should be nil")

	assert.Equal(t, 1, len(emitter.outlets))
	assert.Equal(t, 1, len(dummy1.inlets))
	assert.Equal(t, 1, len(dummy1.outlets))
	assert.Equal(t, 1, len(dummy2.inlets))
	assert.Equal(t, 1, len(dummy2.outlets))
	assert.Equal(t, 1, len(dummy3.inlets))
	assert.Equal(t, 1, len(dummy3.outlets))
	assert.Equal(t, 1, len(sink.inlets))

	pipeline.Run()

	emitter.Emit("Key A", &Vals{})
	emitter.Emit("Key B", &Vals{})
	time.Sleep(time.Millisecond * 20)

	pipeline.Stop()
}

func TestAddDoubleConnect(t *testing.T) {
	emitter := NewEmitterBase("test-emitter", nil)
	sink := NewNullSink("test-sink")
	pipeline := NewPipeline(emitter)

	// Reuse the same wire for the same pair of bolts
	wire, err := pipeline.Plug(emitter, sink)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.PlugWith(emitter, sink, wire)
	assert.NotNil(t, err, "Should Not nil")

	assert.Equal(t, 1, len(emitter.outlets))
	assert.Equal(t, 1, len(sink.inlets))
}
