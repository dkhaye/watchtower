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
	// ErrPayloadTooLarge identifies payloads that exceed MaxPayloadBytes.
	ErrPayloadTooLarge = errors.New("event payload too large")

	// ErrMalformedJSON identifies payloads that are not syntactically valid JSON.
	ErrMalformedJSON = errors.New("malformed event JSON")

	// ErrInvalidEvent identifies valid JSON that does not satisfy the event contract.
	ErrInvalidEvent = errors.New("invalid event")
)

// MaxPayloadBytes is the largest serialized event accepted by Decode.
const MaxPayloadBytes = 1 << 20

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
	if len(payload) > MaxPayloadBytes {
		return Record{}, ErrPayloadTooLarge
	}
	if !json.Valid(payload) {
		return Record{}, ErrMalformedJSON
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(payload, &object); err != nil {
		return Record{}, fmt.Errorf("%w: event must be an object", ErrInvalidEvent)
	}
	if object == nil {
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

	if !hasStrictRFC3339Syntax(timestampText) {
		return Record{}, fmt.Errorf("%w: timestamp must use RFC 3339", ErrInvalidEvent)
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

func hasStrictRFC3339Syntax(value string) bool {
	if len(value) < len("0000-00-00T00:00:00Z") {
		return false
	}

	for i := range 19 {
		switch i {
		case 4, 7:
			if value[i] != '-' {
				return false
			}
		case 10:
			if value[i] != 'T' {
				return false
			}
		case 13, 16:
			if value[i] != ':' {
				return false
			}
		default:
			if value[i] < '0' || value[i] > '9' {
				return false
			}
		}
	}

	zoneStart := 19
	if value[zoneStart] == '.' {
		zoneStart++
		fractionStart := zoneStart
		for zoneStart < len(value) && value[zoneStart] >= '0' && value[zoneStart] <= '9' {
			zoneStart++
		}
		if zoneStart == fractionStart {
			return false
		}
	}

	if zoneStart == len(value)-1 && value[zoneStart] == 'Z' {
		return true
	}
	if len(value)-zoneStart != len("+00:00") {
		return false
	}

	zone := value[zoneStart:]
	if (zone[0] != '+' && zone[0] != '-') || zone[3] != ':' {
		return false
	}
	if zone[1] < '0' || zone[1] > '9' || zone[2] < '0' || zone[2] > '9' ||
		zone[4] < '0' || zone[4] > '9' || zone[5] < '0' || zone[5] > '9' {
		return false
	}

	offsetHour := 10*int(zone[1]-'0') + int(zone[2]-'0')
	offsetMinute := 10*int(zone[4]-'0') + int(zone[5]-'0')
	return offsetHour <= 23 && offsetMinute <= 59
}
