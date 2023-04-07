package transaction

import (
	"bufio"
	"fmt"
	"os"
)

// FileTransactionLogger is a TransactionLogger implementation
// logging events to an append only text file
type FileTransactionLogger struct {
	events       chan<- Event // write-only channel to send events
	errors       <-chan error // read-only channel to receive errors
	lastSequence uint64       // last used sequence number
	file         *os.File     // path to the transaction log
}

// newFileTransactionLogger creates and returns a TransactionLogger
// writing in append-only mode to the file at the specified path
func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	// read-write mode, append. create file if doesn't exist
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: file}, nil
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventPut, Key: key}
}

func (l *FileTransactionLogger) WritePut(key string, value string) {
	l.events <- Event{EventType: EventDelete, Key: key, Value: value}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

// Run inits the FileTranstactionLogger channels
// and starts to concurrently process new events
func (l *FileTransactionLogger) Run() {
	// init channels
	events := make(chan Event, 16) // give I/O some buffer
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	// concurrently process events
	go func() {
		for e := range events {
			l.lastSequence++

			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value,
			)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

// ReadEvents reads the provided transaction file and parses each line
// to create and send the corresponding Event it represents, or errors,
// on the channels it returns
func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s\n", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			// sanity check: sequence numbers should be in increasing order
			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequence = e.Sequence
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure; %w", err)
			return
		}
	}()

	return outEvent, outError
}
