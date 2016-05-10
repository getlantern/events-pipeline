package processors

import (
	"container/list"
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
	timeout time.Duration
	maxEvs  uint64
}

type Slicer struct {
	*events.ProcessorBase

	directives map[events.Key]SlicerDirective
	filtered   map[string]*events.Event
	unfiltered *list.List
	options    *SlicerOptions
}

func NewSlicer(id string, opts *SlicerOptions, ds ...SlicerDirective) *Slicer {

	go func() {
		// TODO:
	}()

	return &Slicer{
		ProcessorBase: events.NewProcessorBase(id, nil),
		filtered:      make(map[string]*events.Event),
		unfiltered:    list.New(),
		directives:    ds,
		options:       opts,
	}
}

func (s *Slicer) Receive(evt *events.Event) error {
	log.Tracef("SLICER ID %v PROCESSED event: %v with: %v", a.ID(), evt.Key, evt.Vals)

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
		l.PushBack(evt)
	}

	//return a.ProcessorBase.Send(evt)
}
