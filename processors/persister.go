package processors

import (
	"bytes"
	"encoding/gob"
	"io"
	"math"
	"os"
	"sync"
	"sync/atomic"

	events "github.com/getlantern/events-pipeline"
)

type PersisterOptions struct {
	MaxBufferSize uint64
	MaxEvents     uint32
	PersistPath   string
}

type Persister struct {
	*events.ProcessorBase
	options *PersisterOptions

	numEvents   uint32
	journalFile *os.File
	persistWg   sync.WaitGroup
	b           *bytes.Buffer
	bMtx        sync.RWMutex
	e           *gob.Encoder
	d           *gob.Decoder
}

func NewPersister(id string, opts *PersisterOptions) *Persister {
	if opts.MaxBufferSize == 0 {
		opts.MaxBufferSize = math.MaxUint64
	}

	if opts.MaxEvents == 0 {
		opts.MaxEvents = math.MaxInt32
	}

	if opts.PersistPath == "" {
		panic("PersisterOptions MUST include PersistPath")
	}

	b := new(bytes.Buffer)
	p := &Persister{
		options: opts,
		b:       b,
		e:       gob.NewEncoder(b),
		d:       gob.NewDecoder(b),
	}
	commitMarker := func(e *events.Event) error { p.markCommit(); return nil }
	p.ProcessorBase = events.NewProcessorBase(id, commitMarker)

	return p
}

func (p *Persister) Receive(evt *events.Event) error {
	log.Tracef("Persister ID %v processed event: %v with: %v", p.ID(), evt.Key, evt.Vals)

	// Handle the SystemEvent signals
	if evt.Key == "" {
		if _, ok := evt.Vals[string(events.SystemEventInit)]; ok {
			log.Debugf("Initializing Persister")

			// First try to recover events is necessary
			var err error
			p.journalFile, err = os.OpenFile(p.options.PersistPath, os.O_RDONLY, 0666)
			if err != nil {
				if !os.IsNotExist(err) {
					log.Errorf("Error opening events recovery file: %v", err)
				}
			} else {
				recovered, errRvr := p.recoverEvents()
				if errRvr != nil {
					log.Errorf("Error processing recovery file: %v", err)
				} else {
					for _, ev := range recovered {
						p.Receive(&ev)
					}
				}
				p.journalFile.Close()
			}

			p.journalFile, err = os.OpenFile(
				p.options.PersistPath,
				os.O_APPEND|os.O_WRONLY|os.O_CREATE,
				0666,
			)
			if err != nil {
				log.Errorf("Error opening or creating event recovery file: %v", err)
			}
		} else if _, ok := evt.Vals[string(events.SystemEventStop)]; ok {
			p.persistWg.Wait()
			p.journalFile.Close()
		}
		return nil
	}

	err := p.ProcessorBase.Receive(evt)
	if err != nil {
		return err
	}

	if err := p.persistEvent(evt); err != nil {
		log.Errorf("Error saving event to recovery file: %v", err)
	}

	return p.ProcessorBase.Send(evt)
}

func (p *Persister) markCommit() {
	// TODO
}

func (p *Persister) persistEvent(evt *events.Event) error {
	p.persistWg.Add(1)

	p.bMtx.Lock()
	err := p.e.Encode(evt)
	if err != nil {
		log.Errorf("Error encoding event to gob: %v", err)
	}
	p.bMtx.Unlock()

	atomic.AddUint32(&p.numEvents, 1)

	p.bMtx.RLock()

	//
	// TODO:
	// - write in chunks (which should be options)
	// - rotating files
	// Keeping sizes reasonable via rotation or similar
	//

	_, err = p.journalFile.Write(p.b.Bytes())
	if err != nil {
		p.bMtx.RUnlock()
		return err

	}
	p.bMtx.RUnlock()

	// We will force a commit mark if any of the limits is reached.  This ensures
	// that in case of failure we don't end up recovering too many events on the next run
	if p.numEvents >= p.options.MaxEvents || p.b.Len() >= int(p.options.MaxBufferSize) {
		p.markCommit()
	}
	p.persistWg.Done()

	return nil
}

func (p *Persister) recoverEvents() ([]events.Event, error) {
	//
	// TODO:
	// Scan the recovery file for the latest commit mark and replay events from there
	//

	p.bMtx.Lock()
	defer p.bMtx.Unlock()

	_, err := io.Copy(p.b, p.journalFile)
	if err != nil {
		return nil, err
	}
	var recovered []events.Event
	var evt events.Event
	for {
		err := p.d.Decode(&evt)

		log.Tracef("Error in decoder: %v", err)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		recovered = append(recovered, evt)
	}

	p.b.Reset()
	return recovered, nil
}
