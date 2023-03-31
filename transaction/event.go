package transaction

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence  uint64    // unique record ID
	EventType EventType // action taken
	Key       string    // key affected by this transaction
	Value     string    // value of a PUT transaction
}
