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

func NewIdentityProcessor(id string) *IdentityProcessor {
	return &IdentityProcessor{
		ProcessorBase: NewProcessorBase(id),
	}
}

func (p *IdentityProcessor) Receive(evt *Event) error {
	log.Tracef("PROCESSOR ID %v received event: %v with: %v", p.ID(), evt.Key, evt.Vals)
	err := p.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	return p.ProcessorBase.Send(evt)
}

func (p *IdentityProcessor) Feedback(evt *Event) error {
	log.Tracef("PROCESSOR ID %v received FEEDBACK of: %v with: %v", p.ID(), evt.Key, evt.Vals)
	return p.ProcessorBase.Feedback(evt)
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

	emitter := NewEmitterBase("test-emitter")
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

func TestIdentityProcessor(t *testing.T) {
	emitter := NewEmitterBase("test-emitter")
	sink := NewNullSink("test-sink")
	dummy := NewIdentityProcessor("test-processor")
	pipeline := NewPipeline(emitter)
	_, err := pipeline.Plug(emitter, dummy)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(dummy, sink)
	assert.Nil(t, err, "Should be nil")

	assert.Equal(t, 1, len(emitter.outlets))
	assert.Equal(t, 1, len(dummy.inlets))
	assert.Equal(t, 1, len(dummy.outlets))
	assert.Equal(t, 1, len(sink.inlets))

	pipeline.Run()

	emitter.Emit("Key A", &Vals{})
	emitter.Emit("Key B", &Vals{})
	time.Sleep(time.Millisecond * 20)

	pipeline.Stop()
}

/*

func TestProcessorChain(t *testing.T) {
	emitter := NewEmitterBase("test-emitter")
	sink := NewNullSink("test-sink")
	dummy1 := NewIdentityProcessor("test-processor A")
	dummy2 := NewIdentityProcessor("test-processor B")
	dummy3 := NewIdentityProcessor("test-processor C")
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
	emitter := NewEmitterBase("test-emitter")
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
*/
