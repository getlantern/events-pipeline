package processors

import (
	"container/list"
	"math"
	"sync"
	"sync/atomic"
	"time"

	events "github.com/getlantern/events-pipeline"
)

const (
	KeepLast   = 1
	KeepFirst  = 2
	KeepRandom = 3 // TODO
)

type DirectiveType int

type SlicerDirective struct {
	Key   events.Key
	dtype DirectiveType
}

type SlicerOptions struct {
	timeout   time.Duration
	maxEvents uint64
}

type Slicer struct {
	*events.ProcessorBase
	directives map[events.Key]SlicerDirective
	filtered   map[events.Key]*events.Event
	unfiltered *list.List
	options    *SlicerOptions

	unfMtx     sync.Mutex
	forceFlush chan struct{}
	numEvs     uint64
}

func NewSlicer(id string, opts *SlicerOptions, ds ...SlicerDirective) *Slicer {
	if opts.maxEvents == 0 {
		opts.maxEvents = math.MaxUint64
	}

	dsmap := make(map[events.Key]SlicerDirective)
	for _, d := range ds {
		dsmap[d.Key] = d
	}

	slicer := &Slicer{
		ProcessorBase: events.NewProcessorBase(id, nil),
		filtered:      make(map[events.Key]*events.Event),
		unfiltered:    list.New(),
		directives:    dsmap,
		options:       opts,
	}

	go func() {
		var timer *time.Timer
		if opts.timeout != 0 {
			timer = time.NewTimer(time.Second * opts.timeout)
		} else {
			timer = time.NewTimer(math.MaxInt64)
		}

		flush := func() {
			slicer.unfMtx.Lock()
			for el := slicer.unfiltered.Front(); el != nil; el = el.Next() {
				err := slicer.ProcessorBase.Send(el.Value.(*events.Event))
				if err != nil {
					log.Errorf("Error sending event")
				}

			}
			slicer.unfiltered = list.New()
			slicer.unfMtx.Unlock()
			// TODO: process filtered
		}

		for {
			select {
			case <-timer.C:
			case <-slicer.forceFlush:
				flush()
			}
		}
	}()

	return slicer
}

func (s *Slicer) Receive(evt *events.Event) error {
	log.Tracef("SLICER ID %v PROCESSED event: %v with: %v", s.ID(), evt.Key, evt.Vals)

	err := s.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	if d, ok := s.directives[evt.Key]; ok {
		switch d.dtype {
		case KeepLast:
			s.filtered[evt.Key] = evt
		case KeepFirst:
			if _, ok := s.filtered[evt.Key]; !ok {
				s.filtered[evt.Key] = evt
			}
		case KeepRandom:
			// TODO: handle the statiscical correctness of this case
		}
	} else {
		s.unfiltered.PushBack(evt)
	}
	if atomic.AddUint64(&s.numEvs, 1) >= s.options.maxEvents {
		s.forceFlush <- struct{}{}
	}

	return nil
}
