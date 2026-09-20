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
	"unicode/utf8"
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
	if !utf8.Valid(payload) || !json.Valid(payload) {
		return Record{}, ErrMalformedJSON
	}
	if !hasValidUnicodeEscapes(payload) {
		return Record{}, fmt.Errorf("%w: event must contain valid Unicode", ErrInvalidEvent)
	}
	if !hasUniqueObjectMemberNames(payload) {
		return Record{}, fmt.Errorf("%w: duplicate object member name", ErrInvalidEvent)
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

	timestamp, err := parseEventTimestamp(timestampText)
	if err != nil {
		return Record{}, fmt.Errorf("%w: timestamp must use the Watchtower timestamp profile", ErrInvalidEvent)
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

func hasUniqueObjectMemberNames(payload []byte) bool {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()
	return nextJSONValueHasUniqueObjectMemberNames(decoder)
}

func nextJSONValueHasUniqueObjectMemberNames(decoder *json.Decoder) bool {
	token, err := decoder.Token()
	if err != nil {
		return false
	}

	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return true
	}

	switch delimiter {
	case '{':
		members := make(map[string]struct{})
		for decoder.More() {
			nameToken, err := decoder.Token()
			if err != nil {
				return false
			}
			name, ok := nameToken.(string)
			if !ok {
				return false
			}
			if _, duplicate := members[name]; duplicate {
				return false
			}
			members[name] = struct{}{}

			if !nextJSONValueHasUniqueObjectMemberNames(decoder) {
				return false
			}
		}
		end, err := decoder.Token()
		return err == nil && end == json.Delim('}')
	case '[':
		for decoder.More() {
			if !nextJSONValueHasUniqueObjectMemberNames(decoder) {
				return false
			}
		}
		end, err := decoder.Token()
		return err == nil && end == json.Delim(']')
	default:
		return false
	}
}

// hasValidUnicodeEscapes rejects lone UTF-16 surrogate escapes anywhere in a
// serialized JSON value before encoding/json replaces them with U+FFFD. JSON
// strings may encode scalar values above U+FFFF as a
// high-surrogate/low-surrogate pair, but an unpaired surrogate is not a Unicode
// scalar value.
func hasValidUnicodeEscapes(raw []byte) bool {
	for i := 1; i < len(raw)-1; {
		if raw[i] != '\\' {
			i++
			continue
		}
		if i+1 >= len(raw)-1 {
			return false
		}
		if raw[i+1] != 'u' {
			i += 2
			continue
		}

		codePoint, ok := decodeHexEscape(raw, i)
		if !ok {
			return false
		}
		i += len(`\u0000`)

		switch {
		case codePoint >= 0xD800 && codePoint <= 0xDBFF:
			lowSurrogate, ok := decodeHexEscape(raw, i)
			if !ok || lowSurrogate < 0xDC00 || lowSurrogate > 0xDFFF {
				return false
			}
			i += len(`\u0000`)
		case codePoint >= 0xDC00 && codePoint <= 0xDFFF:
			return false
		}
	}

	return true
}

func decodeHexEscape(raw []byte, start int) (uint16, bool) {
	if start+len(`\u0000`) > len(raw)-1 || raw[start] != '\\' || raw[start+1] != 'u' {
		return 0, false
	}

	var value uint16
	for _, digit := range raw[start+2 : start+6] {
		value <<= 4
		switch {
		case digit >= '0' && digit <= '9':
			value += uint16(digit - '0')
		case digit >= 'a' && digit <= 'f':
			value += uint16(digit-'a') + 10
		case digit >= 'A' && digit <= 'F':
			value += uint16(digit-'A') + 10
		default:
			return 0, false
		}
	}

	return value, true
}

func parseEventTimestamp(value string) (time.Time, error) {
	if !hasEventTimestampSyntax(value) {
		return time.Time{}, errors.New("invalid event timestamp")
	}

	normalized := []byte(value)
	normalized[10] = 'T'
	if normalized[len(normalized)-1] == 'z' {
		normalized[len(normalized)-1] = 'Z'
	}

	return time.Parse(time.RFC3339Nano, string(normalized))
}

// hasEventTimestampSyntax validates the RFC 3339 profile defined by ADR-0007.
// The profile is intentionally limited to the precision and civil-time values
// that time.Time can represent without normalization: seconds 00-59 and at
// most nine fractional digits.
func hasEventTimestampSyntax(value string) bool {
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
			if value[i] != 'T' && value[i] != 't' {
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

	year := decimalDigits(value[0:4])
	month := decimalDigits(value[5:7])
	day := decimalDigits(value[8:10])
	hour := decimalDigits(value[11:13])
	minute := decimalDigits(value[14:16])
	second := decimalDigits(value[17:19])
	if month < 1 || month > 12 || day < 1 || day > daysInMonth(year, month) ||
		hour > 23 || minute > 59 || second > 59 {
		return false
	}

	zoneStart := 19
	if value[zoneStart] == '.' {
		zoneStart++
		fractionStart := zoneStart
		for zoneStart < len(value) && value[zoneStart] >= '0' && value[zoneStart] <= '9' {
			zoneStart++
		}
		if fractionDigits := zoneStart - fractionStart; fractionDigits < 1 || fractionDigits > 9 {
			return false
		}
	}

	if zoneStart == len(value)-1 && (value[zoneStart] == 'Z' || value[zoneStart] == 'z') {
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

	offsetHour := decimalDigits(zone[1:3])
	offsetMinute := decimalDigits(zone[4:6])
	return offsetHour <= 23 && offsetMinute <= 59
}

func decimalDigits(value string) int {
	result := 0
	for i := range len(value) {
		result = result*10 + int(value[i]-'0')
	}
	return result
}

func daysInMonth(year, month int) int {
	switch month {
	case 4, 6, 9, 11:
		return 30
	case 2:
		if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
			return 29
		}
		return 28
	default:
		return 31
	}
}
