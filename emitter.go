package events

type Emitter interface {
	Emit(Key, *Vals) error
}
