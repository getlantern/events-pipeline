package processors

import (
	"os"
	"testing"
	"time"

	"github.com/getlantern/testify/assert"

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

// Override just to add logging
func (s *NullSink) Receive(evt *events.Event) error {
	log.Tracef("SINK ID %v received event: %v with: %v", s.ID(), evt.Key, evt.Vals)
	return s.SinkBase.Receive(evt)
}

// Callback Sink
type CallbackSink struct {
	*events.SinkBase
	callback func(e *events.Event)
}

func NewCallbackSink(id string, cb func(e *events.Event)) *CallbackSink {
	return &CallbackSink{
		SinkBase: events.NewSinkBase(id),
		callback: cb,
	}
}

func (c *CallbackSink) Receive(e *events.Event) error {
	// Do not perform the callback on system events
	if e.Key != "" {
		c.callback(e)
	}
	return c.SinkBase.Receive(e)
}

func TestIdentityProcessor(t *testing.T) {
	emitter := events.NewEmitterBase("test-emitter", nil)
	sink := NewNullSink("test-sink")
	dummy := NewIdentityProcessor("test-processor")
	pipeline := events.NewPipeline(emitter)
	_, err := pipeline.Plug(emitter, dummy)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(dummy, sink)
	assert.Nil(t, err, "Should be nil")

	pipeline.Run()

	emitter.Emit("Key A", &events.Vals{})
	emitter.Emit("Key B", &events.Vals{})
	time.Sleep(time.Millisecond * 20)

	pipeline.Stop()
}

func TestAggregator(t *testing.T) {
	evs := make(chan *events.Event, 3)

	emitter := events.NewEmitterBase("test-emitter", nil)
	sink := NewCallbackSink("test-sink", func(e *events.Event) {
		log.Tracef("Entering callback with event %v", e)
		evs <- e
	})

	aggregator := NewAggregator(
		"test-aggregator",
		AggregationDirective{"Karma", "level", AggregatorIntRunningSum, RunningSumIdentity},
		AggregationDirective{"Happiness", "level", AggregatorFloat64MovingAverage, MovingAverageIdentity},
	)

	pipeline := events.NewPipeline(emitter)

	_, err := pipeline.Plug(emitter, aggregator)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(aggregator, sink)
	assert.Nil(t, err, "Should be nil")

	pipeline.Run()

	// Test Running Sum
	emitter.Emit("Karma", &events.Vals{"level": 20})
	e := <-evs
	assert.Equal(t, 20, e.Vals["level"], "Should hold this value")

	emitter.Emit("Karma", &events.Vals{"level": 20})
	e = <-evs
	assert.Equal(t, 40, e.Vals["level"], "Should hold this value")

	emitter.Emit("Karma", &events.Vals{"level": 20})
	e = <-evs
	assert.Equal(t, 60, e.Vals["level"], "Should hold this value")

	// Test moving average
	// With these values, we shouldn't have floating point errors, but i
	emitter.Emit("Happiness", &events.Vals{"level": 250.5})
	e = <-evs
	assert.Equal(t, 250.5, e.Vals["level"], "Should hold this value")

	emitter.Emit("Happiness", &events.Vals{"level": 0.5})
	e = <-evs
	assert.Equal(t, 125.5, e.Vals["level"], "Should hold this value")

	emitter.Emit("Happiness", &events.Vals{"level": 300.0})
	emitter.Emit("Happiness", &events.Vals{"level": 400.0})
	e = <-evs
	e = <-evs
	assert.Equal(t, 237.75, e.Vals["level"], "Should hold this value")

	pipeline.Stop()
}

func TestCondenser(t *testing.T) {
	evs := make(chan *events.Event, 3)

	emitter := events.NewEmitterBase("test-emitter", nil)

	sink := NewCallbackSink("test-sink", func(e *events.Event) {
		log.Tracef("Entering callback with event %v", e)
		evs <- e
	})

	condenser := NewCondenser(
		"test-condenser",
		&CondenserOptions{Timeout: time.Minute, MaxEvents: 2},
	)

	pipeline := events.NewPipeline(emitter)
	_, err := pipeline.Plug(emitter, condenser)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(condenser, sink)
	assert.Nil(t, err, "Should be nil")

	pipeline.Run()

	emitter.Emit("Empathy", &events.Vals{})
	assert.Equal(t, 0, len(evs), "Event shouldn't have reached the sink")

	emitter.Emit("Empathy", &events.Vals{})
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 2, len(evs), "Both events should have reached the sink")

	pipeline.Stop()
}

func TestKeyRateLimiter(t *testing.T) {
	evs := make(chan *events.Event, 3)

	emitter := events.NewEmitterBase("test-emitter", nil)
	sink := NewCallbackSink("test-sink", func(e *events.Event) {
		log.Tracef("Entering callback with event %v", e)
		evs <- e
	})

	ratelimiter := NewKeyRateLimiter(
		"test-ratelimiter",
		&KeyRateLimiterOptions{Interval: time.Minute, MaxPerInterval: 2},
	)

	pipeline := events.NewPipeline(emitter)
	_, err := pipeline.Plug(emitter, ratelimiter)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(ratelimiter, sink)
	assert.Nil(t, err, "Should be nil")

	pipeline.Run()

	emitter.Emit("Wisdom", &events.Vals{})
	emitter.Emit("Wisdom", &events.Vals{})
	emitter.Emit("Wisdom", &events.Vals{})

	time.Sleep(20 * time.Millisecond)

	assert.Equal(t, 2, len(evs), "Only two events should have arrived to the sink")

	pipeline.Stop()
}

func TestPersister(t *testing.T) {
	persistPath := "test-persister"
	defer func() {
		if _, err := os.Stat(persistPath); err == nil {
			err := os.Remove(persistPath)
			if err != nil {
				log.Errorf("Error removing Persister recovery file: %v", err)
				return
			}
		}
	}()

	emitter := events.NewEmitterBase("test-emitter", nil)
	sink := NewNullSink("test-sink")
	persister := NewPersister(
		"test-processor",
		&PersisterOptions{
			MaxEvents:   2,
			PersistPath: persistPath,
		})
	pipeline := events.NewPipeline(emitter)
	_, err := pipeline.Plug(emitter, persister)
	assert.Nil(t, err, "Should be nil")
	_, err = pipeline.Plug(persister, sink)
	assert.Nil(t, err, "Should be nil")

	pipeline.Run()

	evt := events.NewEvent("Colors", &events.Vals{"Beauty": "Imperfection"})

	if err := persister.persistEvent(evt); err != nil {
		assert.Nil(t, err, "Error should be nil")
	}

	// Force opening with both read and write permissions, defying the logic of the Persister
	persister.journalFile, err = os.OpenFile(
		persister.options.PersistPath,
		os.O_RDWR|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Errorf("Error opening or creating event recovery file: %v", err)
	}

	persister.b.Reset()

	evts, err := persister.recoverEvents()
	if err != nil {
		t.Fatalf("Error recovering events: %v", err)
	}

	if assert.Equal(t, 1, len(evts), "One event should have been recovered") {
		assert.Equal(t, *evt, evts[0], "The recovered event should be the last one emitted")
	}

	time.Sleep(100 * time.Millisecond)

	pipeline.Stop()
}
