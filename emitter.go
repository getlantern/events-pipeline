package events

type Emitter interface {
	Emit(*Key, *Vals)
}

type DefaultEmitter struct {
	Sender
}

func NewEmitter() *DefaultEmitter {
	return &DefaultEmitter{}
}

func (e *DefaultEmitter) Emit(k *Key, v *Vals) error {
	return e.Sender.Send(MakeEvent(k, v))
}

func (e *DefaultEmitter) Link(r *Receiver) {
}
