package processors

import (
	"testing"
	"time"

	"github.com/getlantern/testify/assert"

	events "github.com/getlantern/events-pipeline"
)

func TestAggregator(t *testing.T) {
	emitter := events.NewEmitter("test-emitter")
	sink := events.NewNullSink("test-sink")
	aggregator := NewAggregator("test-aggregator", nil)

	pipeline := events.NewPipeline(emitter)

	_, err := pipeline.Plug(emitter, aggregator)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(aggregator, sink)
	assert.Nil(t, err, "Should be nil")

	pipeline.Run()

	emitter.Emit("Key A", &events.Vals{})
	emitter.Emit("Key B", &events.Vals{})
	time.Sleep(time.Millisecond * 20)

	pipeline.Stop()
}
