package events

type DummyProcessor struct {
	ProcessorBase
}

func NewDummyProcessor() *DummyProcessor {
	return &DummyProcessor{}
}

func (p *DummyProcessor) LinkOut(w *Wire) {
	p.SenderBase.LinkOutlet(w)
}

func (p *DummyProcessor) LinkInlet(w *Wire) {
	p.ReceiverBase.LinkInlet(w)
}

func (p *DummyProcessor) Receive(evt *Event) error {
	return nil
}

func (p *DummyProcessor) Send(evt *Event) error {
	return p.ProcessorBase.Send(evt)
}
