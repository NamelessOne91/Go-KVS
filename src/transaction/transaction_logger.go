package transaction

import (
	"fmt"

	"github.com/NamelessOne91/Go-KVS/store"
)

// TransactionLogger will be implemented by services which can provide
// the ability to persist an ordered list of mutating events executed by the data store
type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key string, value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

var Logger TransactionLogger

// InitTransactionLogger creates a new TransactionLogger,
// process previously logged events to reproduce the last store state
// and starts processing new events
func InitTransactionLogger() error {
	var err error

	Logger, err = NewPostgresTransactionLogger(PostgresDBParams{
		host:     "localhost",
		name:     "kvs",
		user:     "test",
		password: "test!1",
	})
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := Logger.ReadEvents()
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = store.Delete(e.Key)
			case EventPut:
				err = store.Put(e.Key, e.Value)
			}

		}
	}

	Logger.Run()
	return err
}
