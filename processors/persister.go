package processors

import (
	"bytes"
	"encoding/gob"
	"math"
	"os"

	events "github.com/getlantern/events-pipeline"
)

type PersisterOptions struct {
	MaxBufferSize int
	MaxEvents     int
	PersistPath   string
}

type Persister struct {
	*events.ProcessorBase
	options *PersisterOptions

	numEvents int
	pf        *os.File
	b         *bytes.Buffer
	e         *gob.Encoder
	d         *gob.Decoder
}

func NewPersister(id string, opts *PersisterOptions) *Persister {
	if opts.MaxBufferSize == 0 {
		opts.MaxBufferSize = math.MaxInt64
	}

	if opts.MaxEvents == 0 {
		opts.MaxEvents = math.MaxInt64
	}

	if opts.PersistPath == "" {
		panic("PersisterOptions MUST include PersistPath")
	}

	b := new(bytes.Buffer)
	p := &Persister{
		options: opts,
		b:       b,
		e:       gob.NewEncoder(b),
	}
	commitMarker := func(e *events.Event) error { p.markCommit(); return nil }
	p.ProcessorBase = events.NewProcessorBase(id, commitMarker)

	return p
}

func (p *Persister) Receive(evt *events.Event) error {
	log.Tracef("PERSISTER ID %v PROCESSED event: %v with: %v", p.ID(), evt.Key, evt.Vals)

	// Handle the SystemEvent signals
	if evt.Key == "" {
		if _, ok := evt.Vals[string(events.SystemEventInit)]; ok {
			var err error
			p.pf, err = os.OpenFile(p.options.PersistPath, os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				log.Errorf("Error opening events recovery file")
			} else {
				p.recoverEvents()
			}
		} else if _, ok := evt.Vals[string(events.SystemEventStop)]; ok {
			p.pf.Close()
		}
		return nil
	}

	err := p.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	p.persistEvent(evt)

	return p.ProcessorBase.Send(evt)
}

func (p *Persister) markCommit() {
	// TODO
}

func (p *Persister) persistEvent(evt *events.Event) {
	err := p.e.Encode(evt)
	if err != nil {
		log.Errorf("Error encoding event to gob: %v", err)
	}
	p.numEvents = p.numEvents + 1

	// TODO: write in chunks (which should be options)
	n, err := p.pf.Write(p.b.Bytes())
	if err != nil {
		log.Errorf("Error saving %v bytes to event recovery file: %v", n, err)
	}

	// We will force a commit mark if any of the limits is reached.  This ensures
	// that in case of failure we don't end up recovering too many events on the next run
	if p.numEvents >= p.options.MaxEvents || p.b.Len() >= p.options.MaxBufferSize {
		p.markCommit()
	}

	log.Tracef("BUFFER!: %v")
}

func (p *Persister) recoverEvents() {
	// TODO
	// Scan the recovery file for the latest commit mark and replay events from there
}
