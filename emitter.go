package events

type Emitter interface {
	Sender
	Emit(*Key, *Vals)
}

type DefaultEmitter struct {
}

func NewEmitter() *DefaultEmitter {
	return &DefaultEmitter{}
}

func (e *DefaultEmitter) Emit(k *Key, v *Vals) {
}

func (e *DefaultEmitter) Send(ev *Event) error {
	return nil
}

func (e *DefaultEmitter) Link(r *Receiver) {
}
