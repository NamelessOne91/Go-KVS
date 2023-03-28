package store

import (
	"errors"
	"testing"
)

// basic, quite unnecessary, tests
func TestStore(t *testing.T) {

	_, err := Get("empty-start")
	if !errors.Is(err, ErrorNoSuchKey) {
		t.Errorf("Expected error: %s - got %v", ErrorNoSuchKey, err)
	}

	Put("1", "1")
	Put("2", "2")
	Put("3", "3")

	if len(store) != 3 {
		t.Errorf("Expected store to contains 3 elements - got %d", len(store))
	}

	Put("1", "0")
	if len(store) != 3 {
		t.Errorf("Expected store to contains 3 elements - got %d", len(store))
	}
	v, err := Get("1")
	if err != nil || v != "0" {
		t.Errorf("Expected Get(\"1\") to return \"0\" - got %s", v)
	}

	_ = Delete("1")
	_, err = Get("1")
	if !errors.Is(err, ErrorNoSuchKey) {
		t.Errorf("Expected error: %s - got %v", ErrorNoSuchKey, err)
	}
}
