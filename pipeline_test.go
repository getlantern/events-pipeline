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
