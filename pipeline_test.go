package events

import (
	"testing"
	"time"

	"fmt"

	"github.com/getlantern/testify/assert"
)

func TestTrivialPipeline(t *testing.T) {
	emitter := NewEmitter()
	sink := NewDummySink()
	pipeline := NewPipeline(emitter)
	pipeline.Plug(emitter, sink)

	assert.Equal(t, 1, len(emitter.outlets))
	assert.Equal(t, 1, len(sink.inlets))

	pipeline.Run()

	fmt.Println("Emit k/v")

	emitter.Emit("Key A", &Vals{})
	emitter.Emit("Key B", &Vals{})
	time.Sleep(time.Millisecond * 20)

	pipeline.Stop()
}

func TestDummyProcessor(t *testing.T) {
	emitter := NewEmitter()
	sink := NewDummySink()
	dummy := NewDummyProcessor()
	pipeline := NewPipeline(emitter)
	pipeline.Plug(emitter, dummy)
	pipeline.Plug(dummy, sink)

	assert.Equal(t, 1, len(emitter.outlets))
	assert.Equal(t, 1, len(dummy.inlets))
	assert.Equal(t, 1, len(dummy.outlets))
	assert.Equal(t, 1, len(sink.inlets))

	pipeline.Run()

	emitter.Emit("Key A", &Vals{})
	//emitter.Emit("Key B", &Vals{})
	time.Sleep(time.Millisecond * 20)

	pipeline.Stop()
}
