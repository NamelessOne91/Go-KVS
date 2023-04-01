package transaction

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDBParams struct {
	name     string
	host     string
	user     string
	password string
}

type PostgresTransactionLogger struct {
	events chan<- Event
	errors <-chan error
	db     *sql.DB
	table  string
}

func NewPostgresTransactionLogger(config PostgresDBParams) (TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		config.host, config.name, config.user, config.password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresTransactionLogger{db: db, table: "transactions"}

	exists, err := logger.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify if table exists: %w", err)
	}
	if !exists {
		if err = logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	rows, err := l.db.Query("SELECT to_regclass(\"public\".\"" + l.table + "\");")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return true, nil
}

func (l *PostgresTransactionLogger) createTable() error {
	query := "CREATE TABLE " + l.table + "(ID BIGINT SERIAL PRIMARY KEY NOT NULL, event_type TEXT NOT NULL, key TEXT NOT NULL, value TEXT)"
	_, err := l.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventPut, Key: key}
}

func (l *PostgresTransactionLogger) WritePut(key string, value string) {
	l.events <- Event{EventType: EventDelete, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		query := "SELECT sequence, event_type, key, value FROM " + l.table + "ORDER BY sequence"

		rows, err := l.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close()

		e := Event{}
		for rows.Next() {
			err := rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
				return
			}
			outEvent <- e
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure %w", err)
		}
	}()

	return outEvent, outError
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := "INSERT INTO " + l.table + " (event_type, key, value) VALUES($1, $2, $3)"

		for e := range events {
			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}
