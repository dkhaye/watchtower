// Package event defines Watchtower's initial telemetry event contract and the
// boundary that decodes serialized telemetry into that contract.
package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrMalformedJSON identifies payloads that are not syntactically valid JSON.
	ErrMalformedJSON = errors.New("malformed event JSON")

	// ErrInvalidEvent identifies valid JSON that does not satisfy the event contract.
	ErrInvalidEvent = errors.New("invalid event")
)

// Event is the smallest normalized security-telemetry event Watchtower
// currently understands. The contract will evolve as real telemetry sources
// expose requirements that justify additional fields.
type Event struct {
	ID        string
	Timestamp time.Time
	Actor     string
	Action    string
	Target    string
	Source    string
}

// Record keeps a decoded event coupled to the exact serialized payload from
// which it was produced. The raw payload is copied on input and output so
// callers cannot mutate the record after decoding.
type Record struct {
	Event Event
	raw   []byte
}

// Raw returns a copy of the exact serialized payload supplied to Decode.
func (r Record) Raw() []byte {
	return bytes.Clone(r.raw)
}

// Decode validates one serialized JSON event and returns its typed form
// alongside an immutable copy of the original payload. Unknown object fields
// are accepted and remain available in the raw payload.
func Decode(payload []byte) (Record, error) {
	if !json.Valid(payload) {
		return Record{}, ErrMalformedJSON
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(payload, &object); err != nil {
		return Record{}, fmt.Errorf("%w: event must be an object", ErrInvalidEvent)
	}

	var id, timestampText, actor, action, target, source string
	required := []struct {
		name  string
		value *string
	}{
		{name: "id", value: &id},
		{name: "timestamp", value: &timestampText},
		{name: "actor", value: &actor},
		{name: "action", value: &action},
		{name: "target", value: &target},
		{name: "source", value: &source},
	}
	for _, field := range required {
		raw, exists := object[field.name]
		if !exists {
			return Record{}, fmt.Errorf("%w: %s is required", ErrInvalidEvent, field.name)
		}
		if err := json.Unmarshal(raw, field.value); err != nil {
			return Record{}, fmt.Errorf("%w: %s must be a string", ErrInvalidEvent, field.name)
		}
		if strings.TrimSpace(*field.value) == "" {
			return Record{}, fmt.Errorf("%w: %s is required", ErrInvalidEvent, field.name)
		}
	}

	timestamp, err := time.Parse(time.RFC3339Nano, timestampText)
	if err != nil {
		return Record{}, fmt.Errorf("%w: timestamp must use RFC 3339", ErrInvalidEvent)
	}

	return Record{
		Event: Event{
			ID:        id,
			Timestamp: timestamp,
			Actor:     actor,
			Action:    action,
			Target:    target,
			Source:    source,
		},
		raw: bytes.Clone(payload),
	}, nil
}
