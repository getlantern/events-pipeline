package events

type EmitterBase struct {
	SenderBase
}

func NewEmitter() *EmitterBase {
	return &EmitterBase{}
}

func (e *EmitterBase) Emit(k Key, v *Vals) error {
	return e.SenderBase.Send(MakeEvent(k, v))
}

/*
func (e *EmitterBase) Link(r *Wire) {
	e.SenderBase.Link(r)
}

func (e *EmitterBase) Send(evt *Event) error {
	return nil
}
*/
