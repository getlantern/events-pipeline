package events

type Emitter interface {
	Sender
	Emit(*Key, *Vals)
}
