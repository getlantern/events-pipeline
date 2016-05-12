package processors

import (
	"sync"
	"time"

	events "github.com/getlantern/events-pipeline"
)

type RateLimiterOptions struct {
	Interval       time.Duration
	MaxPerInterval int64
}

type RateLimiter struct {
	*events.ProcessorBase
	options *RateLimiterOptions

	sentKeyCount      keyCountMap
	discardedKeyCount keyCountMap
	keysMtx           sync.Mutex
	ticker            *time.Ticker
}

func NewRateLimiter(id string, opts *RateLimiterOptions) *RateLimiter {
	if opts.MaxPerInterval == 0 {
		panic("Limiting the number of events to 0 per time unit makes no sense")
	}
	if opts.Interval == 0 {
		opts.Interval = time.Minute
	}

	return &RateLimiter{
		ProcessorBase:     events.NewProcessorBase(id, nil),
		options:           opts,
		sentKeyCount:      make(keyCountMap),
		discardedKeyCount: make(keyCountMap),
		ticker:            time.NewTicker(opts.Interval),
	}
}

func (r *RateLimiter) Receive(evt *events.Event) error {
	log.Tracef("RATELIMITER ID %v PROCESSED event: %v with: %v", r.ID(), evt.Key, evt.Vals)

	// Handle the SystemEvent signals
	if evt.Key == "" {
		return nil
	}

	err := r.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	r.keysMtx.Lock()
	defer r.keysMtx.Unlock()

	select {
	case <-r.ticker.C:
		log.Tracef("Sent %v events during the last period", len(r.sentKeyCount))
		log.Tracef("Discarded %v events during the last period", len(r.discardedKeyCount))

		r.sentKeyCount = make(keyCountMap)
		r.discardedKeyCount = make(keyCountMap)
		r.sentKeyCount[evt.Key] = 1

		return r.ProcessorBase.Send(evt)
	default:
		v, ok := r.sentKeyCount[evt.Key]
		if !ok {
			v = 0
		}

		if v < r.options.MaxPerInterval {
			r.sentKeyCount[evt.Key] = v + 1
			return r.ProcessorBase.Send(evt)
		} else {
			if vd, vok := r.discardedKeyCount[evt.Key]; vok {
				r.discardedKeyCount[evt.Key] = vd + 1
			} else {
				r.discardedKeyCount[evt.Key] = 1
			}
			return nil
		}
	}
}
