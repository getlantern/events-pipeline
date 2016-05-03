package events

type Key string
type Vals map[string]interface{}

type Event struct {
	Key  Key
	Vals Vals
}
