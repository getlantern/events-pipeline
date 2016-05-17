// A condenser accumulates events until an event happens (timeout or max events)
// Then it sends them in a burst, followed by a SystemEventMark event

package processors

import (
	"container/list"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	events "github.com/getlantern/events-pipeline"
)

const (
	KeepLast   = 1
	KeepFirst  = 2
	KeepRandom = 3
)

type DirectiveType int

type CondenserDirective struct {
	Key   events.Key
	dtype DirectiveType
}

type CondenserOptions struct {
	Timeout   time.Duration
	MaxEvents uint64
}

type directiveMap map[events.Key]CondenserDirective
type filteredMap map[events.Key]*events.Event

type Condenser struct {
	*events.ProcessorBase
	directives directiveMap
	filtered   filteredMap
	unfiltered *list.List
	options    *CondenserOptions

	keyCount keyCountMap
	r        *rand.Rand

	evMtx      sync.Mutex
	forceFlush chan struct{}
	numEvs     uint64
}

func NewCondenser(id string, opts *CondenserOptions, ds ...CondenserDirective) *Condenser {
	if opts.MaxEvents == 0 {
		opts.MaxEvents = math.MaxUint64
	}

	dsmap := make(directiveMap)
	for _, d := range ds {
		dsmap[d.Key] = d
	}

	s := &Condenser{
		ProcessorBase: events.NewProcessorBase(id, nil),
		directives:    dsmap,
		filtered:      make(filteredMap),
		unfiltered:    list.New(),
		options:       opts,
		keyCount:      make(keyCountMap),
		r:             rand.New(rand.NewSource(time.Now().UnixNano())),
		forceFlush:    make(chan struct{}),
	}

	go func() {
		// TODO: handle stop and restart (need to stop ticker)
		var ticker *time.Ticker
		if opts.Timeout != 0 {
			ticker = time.NewTicker(time.Second * opts.Timeout)
		} else {
			ticker = time.NewTicker(math.MaxInt64)
		}

		for {
			select {
			case <-ticker.C:
			case <-s.forceFlush:
				s.flush()
			}
		}
	}()

	return s
}

func (s *Condenser) Receive(evt *events.Event) error {
	log.Tracef("CONDENSER ID %v PROCESSED event: %v with: %v", s.ID(), evt.Key, evt.Vals)

	// Handle the SystemEvent signals
	if evt.Key == "" {
		if _, ok := evt.Vals[string(events.SystemEventStop)]; ok {
			s.forceFlush <- struct{}{}
		}
		return nil
	}

	err := s.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	s.evMtx.Lock()
	if d, ok := s.directives[evt.Key]; ok {
		switch d.dtype {
		case KeepLast:
			s.filtered[evt.Key] = evt
		case KeepFirst:
			if _, ok := s.filtered[evt.Key]; !ok {
				s.filtered[evt.Key] = evt
			}
		case KeepRandom:
			numEvs, ok := s.keyCount[evt.Key]
			if !ok {
				s.keyCount[evt.Key] = 0
				numEvs = 0
			}
			if s.r.Int63n(numEvs+1) == numEvs {
				s.filtered[evt.Key] = evt
			}
			s.keyCount[evt.Key] = numEvs + 1
		}
	} else {
		s.unfiltered.PushBack(evt)
	}
	s.evMtx.Unlock()

	if atomic.AddUint64(&s.numEvs, 1) >= s.options.MaxEvents {
		s.forceFlush <- struct{}{}
	}

	return nil
}

func (s *Condenser) flush() {
	s.evMtx.Lock()
	for el := s.unfiltered.Front(); el != nil; el = el.Next() {
		err := s.ProcessorBase.Send(el.Value.(*events.Event))
		if err != nil {
			log.Errorf("Error sending event")
		}
	}
	s.unfiltered = list.New()

	for _, v := range s.filtered {
		err := s.ProcessorBase.Send(v)
		if err != nil {
			log.Errorf("Error sending event")
		}
	}
	s.filtered = make(filteredMap)
	s.keyCount = make(keyCountMap)
	s.evMtx.Unlock()

	// Send a SystemEventMark
	err := s.ProcessorBase.Send(events.NewEvent("", &events.Vals{string(events.SystemEventMark): nil}))
	if err != nil {
		log.Errorf("Error sending event")
	}
}
