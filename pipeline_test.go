package events

import (
	"testing"
	"time"

	"github.com/getlantern/testify/assert"
)

func TestTrivialPipeline(t *testing.T) {
	emitter := NewEmitter("test-emitter")
	sink := NewDummySink("test-sink")
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

func TestDummyProcessor(t *testing.T) {
	emitter := NewEmitter("test-emitter")
	sink := NewDummySink("test-sink")
	dummy := NewDummyProcessor("test-processor")
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

func TestProcessorChain(t *testing.T) {
	emitter := NewEmitter("test-emitter")
	sink := NewDummySink("test-sink")
	dummy1 := NewDummyProcessor("test-processor A")
	dummy2 := NewDummyProcessor("test-processor B")
	dummy3 := NewDummyProcessor("test-processor C")
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
	emitter := NewEmitter("test-emitter")
	sink := NewDummySink("test-sink")
	pipeline := NewPipeline(emitter)

	// Reuse the same wire for the same pair of bolts
	wire, err := pipeline.Plug(emitter, sink)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.PlugWith(emitter, sink, wire)
	assert.NotNil(t, err, "Should Not nil")

	assert.Equal(t, 1, len(emitter.outlets))
	assert.Equal(t, 1, len(sink.inlets))
}
