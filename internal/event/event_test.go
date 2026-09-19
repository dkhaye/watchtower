package event_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/dkhaye/watchtower/internal/event"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	payload := []byte(`{
  "id": "event-123",
  "timestamp": "2026-09-19T14:30:45.123456789-04:00",
  "actor": "user@example.com",
  "action": "role.assignment.created",
  "target": "/subscriptions/example/resourceGroups/security",
  "source": "azure-activity-log",
  "schema_version": 7,
  "unknown": {"large_number": 1e100000, "nested": [true, {"value": 1}]}
}`)
	wantRaw := bytes.Clone(payload)

	record, err := event.Decode(payload)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	wantTimestamp := time.Date(2026, time.September, 19, 14, 30, 45, 123456789, time.FixedZone("", -4*60*60))
	if record.Event.ID != "event-123" {
		t.Errorf("Decode() ID = %q, want %q", record.Event.ID, "event-123")
	}
	if !record.Event.Timestamp.Equal(wantTimestamp) {
		t.Errorf("Decode() Timestamp = %v, want %v", record.Event.Timestamp, wantTimestamp)
	}
	if record.Event.Actor != "user@example.com" {
		t.Errorf("Decode() Actor = %q, want %q", record.Event.Actor, "user@example.com")
	}
	if record.Event.Action != "role.assignment.created" {
		t.Errorf("Decode() Action = %q, want %q", record.Event.Action, "role.assignment.created")
	}
	if record.Event.Target != "/subscriptions/example/resourceGroups/security" {
		t.Errorf("Decode() Target = %q, want %q", record.Event.Target, "/subscriptions/example/resourceGroups/security")
	}
	if record.Event.Source != "azure-activity-log" {
		t.Errorf("Decode() Source = %q, want %q", record.Event.Source, "azure-activity-log")
	}
	if got := record.Raw(); !bytes.Equal(got, wantRaw) {
		t.Errorf("Decode() Raw = %q, want exact payload %q", got, wantRaw)
	}

	payload[0] = '['
	if got := record.Raw(); !bytes.Equal(got, wantRaw) {
		t.Errorf("Decode() retained Raw = %q after input mutation, want %q", got, wantRaw)
	}

	returnedRaw := record.Raw()
	returnedRaw[0] = '['
	if got := record.Raw(); !bytes.Equal(got, wantRaw) {
		t.Errorf("Decode() retained Raw = %q after output mutation, want %q", got, wantRaw)
	}
}

func TestDecodeRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	_, err := event.Decode([]byte(`{"id":`))
	if !errors.Is(err, event.ErrMalformedJSON) {
		t.Fatalf("Decode() error = %v, want ErrMalformedJSON", err)
	}
}

func TestDecodeRejectsInvalidUTF8(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"id":"x","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`)
	payload[7] = 0xff

	_, err := event.Decode(payload)
	if !errors.Is(err, event.ErrMalformedJSON) {
		t.Fatalf("Decode() error = %v, want ErrMalformedJSON", err)
	}
}

func TestDecodeRejectsUnpairedUnicodeSurrogateEscapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "high surrogate",
			payload: `{"id":"\ud800","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`,
		},
		{
			name:    "different high surrogate",
			payload: `{"id":"\ud801","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`,
		},
		{
			name:    "low surrogate",
			payload: `{"id":"\udc00","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`,
		},
		{
			name:    "high surrogate followed by non-low surrogate",
			payload: `{"id":"\ud800\u0041","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`,
		},
		{
			name:    "unknown string value",
			payload: `{"id":"1","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d","unknown":"\ud800"}`,
		},
		{
			name:    "unknown member name",
			payload: `{"id":"1","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d","\ud800":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := event.Decode([]byte(tt.payload))
			if !errors.Is(err, event.ErrInvalidEvent) {
				t.Fatalf("Decode() error = %v, want ErrInvalidEvent", err)
			}
			if !strings.Contains(err.Error(), "event must contain valid Unicode") {
				t.Errorf("Decode() error = %q, want Unicode context", err)
			}
		})
	}
}

func TestDecodeAcceptsValidUnicodeEscapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   string
		want string
	}{
		{name: "surrogate pair", id: `\ud83d\ude80`, want: "🚀"},
		{name: "replacement character escape", id: `\ufffd`, want: "�"},
		{name: "escaped backslash", id: `\\ud800`, want: `\ud800`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			payload := []byte(`{"id":"` + tt.id + `","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`)
			record, err := event.Decode(payload)
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			if record.Event.ID != tt.want {
				t.Errorf("Decode() ID = %q, want %q", record.Event.ID, tt.want)
			}
		})
	}
}

func TestDecodeRejectsDuplicateObjectMemberNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "required member",
			payload: `{"id":"first","id":"second","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`,
		},
		{
			name:    "escaped equivalent member",
			payload: `{"id":"first","\u0069d":"second","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`,
		},
		{
			name:    "unknown member",
			payload: `{"id":"first","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d","extra":1,"extra":2}`,
		},
		{
			name:    "nested member",
			payload: `{"id":"first","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d","extra":{"value":1,"value":2}}`,
		},
		{
			name:    "member nested in array",
			payload: `{"id":"first","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d","extra":[{"value":1,"value":2}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := event.Decode([]byte(tt.payload))
			if !errors.Is(err, event.ErrInvalidEvent) {
				t.Fatalf("Decode() error = %v, want ErrInvalidEvent", err)
			}
			if !strings.Contains(err.Error(), "duplicate object member name") {
				t.Errorf("Decode() error = %q, want duplicate-member context", err)
			}
			if strings.Contains(err.Error(), "first") || strings.Contains(err.Error(), "extra") {
				t.Errorf("Decode() error %q contains telemetry field data", err)
			}
		})
	}
}

func TestDecodeAcceptsRFC3339CaseVariants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		timestamp string
	}{
		{name: "lowercase time separator", timestamp: "2026-09-19t18:30:45Z"},
		{name: "lowercase UTC designator", timestamp: "2026-09-19T18:30:45z"},
		{name: "both lowercase", timestamp: "2026-09-19t18:30:45.123z"},
		{name: "lowercase time separator with offset", timestamp: "2026-09-19t18:30:45-04:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			payload := []byte(`{"id":"1","timestamp":"` + tt.timestamp + `","actor":"a","action":"b","target":"c","source":"d"}`)
			record, err := event.Decode(payload)
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			if got := record.Raw(); !bytes.Equal(got, payload) {
				t.Errorf("Decode() Raw = %q, want exact payload %q", got, payload)
			}
		})
	}
}

func TestDecodeRejectsOversizedPayloadBeforeParsing(t *testing.T) {
	t.Parallel()

	payload := bytes.Repeat([]byte("x"), event.MaxPayloadBytes+1)
	_, err := event.Decode(payload)
	if !errors.Is(err, event.ErrPayloadTooLarge) {
		t.Fatalf("Decode() error = %v, want ErrPayloadTooLarge", err)
	}
}

func TestDecodeRejectsInvalidEvents(t *testing.T) {
	t.Parallel()

	valid := `{
  "id": "event-123",
  "timestamp": "2026-09-19T18:30:45Z",
  "actor": "user@example.com",
  "action": "role.assignment.created",
  "target": "security",
  "source": "azure-activity-log"
}`

	tests := []struct {
		name    string
		payload string
		field   string
	}{
		{name: "null", payload: `null`, field: "object"},
		{name: "array", payload: `[]`, field: "object"},
		{name: "wrong field type", payload: `{"id":{},"timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`, field: "id must be a string"},
		{name: "wrong field case", payload: strings.Replace(valid, `"id": "event-123"`, `"ID": "event-123"`, 1), field: "id"},
		{name: "missing id", payload: strings.Replace(valid, `"id": "event-123",`, ``, 1), field: "id"},
		{name: "missing timestamp", payload: strings.Replace(valid, `"timestamp": "2026-09-19T18:30:45Z",`, ``, 1), field: "timestamp"},
		{name: "missing actor", payload: strings.Replace(valid, `"actor": "user@example.com",`, ``, 1), field: "actor"},
		{name: "missing action", payload: strings.Replace(valid, `"action": "role.assignment.created",`, ``, 1), field: "action"},
		{name: "missing target", payload: strings.Replace(valid, `"target": "security",`, ``, 1), field: "target"},
		{name: "missing source", payload: strings.Replace(valid, `,
  "source": "azure-activity-log"`, ``, 1), field: "source"},
		{name: "blank value", payload: strings.Replace(valid, `"actor": "user@example.com"`, `"actor": "  "`, 1), field: "actor"},
		{name: "invalid timestamp", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "yesterday", 1), field: "timestamp"},
		{name: "non-digit timestamp", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "202X-09-19T18:30:45Z", 1), field: "timestamp"},
		{name: "invalid date separator", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026/09-19T18:30:45Z", 1), field: "timestamp"},
		{name: "empty timestamp fraction", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45.Z", 1), field: "timestamp"},
		{name: "comma timestamp fraction", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45,1Z", 1), field: "timestamp"},
		{name: "missing timestamp zone", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45.1", 1), field: "timestamp"},
		{name: "invalid timestamp zone", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45X", 1), field: "timestamp"},
		{name: "compact timestamp offset", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45+0100", 1), field: "timestamp"},
		{name: "invalid offset sign", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45*01:00", 1), field: "timestamp"},
		{name: "invalid offset separator", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45+01-00", 1), field: "timestamp"},
		{name: "non-digit offset hour", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45+0x:00", 1), field: "timestamp"},
		{name: "non-digit offset minute", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45+01:0x", 1), field: "timestamp"},
		{name: "invalid offset minute", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45+01:60", 1), field: "timestamp"},
		{name: "invalid offset hour", payload: strings.Replace(valid, "2026-09-19T18:30:45Z", "2026-09-19T18:30:45+24:00", 1), field: "timestamp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := event.Decode([]byte(tt.payload))
			if !errors.Is(err, event.ErrInvalidEvent) {
				t.Fatalf("Decode() error = %v, want ErrInvalidEvent", err)
			}
			if !strings.Contains(err.Error(), tt.field) {
				t.Errorf("Decode() error = %q, want field context %q", err, tt.field)
			}
		})
	}
}

func TestDecodeDoesNotIncludePayloadInErrors(t *testing.T) {
	t.Parallel()

	const payloadMarker = "payload-value-must-not-appear"
	payloads := [][]byte{
		[]byte(`{"id":"` + payloadMarker + `"}`),
		[]byte(`{"id":"1","timestamp":"` + payloadMarker + `","actor":"a","action":"b","target":"c","source":"d"}`),
	}

	for _, payload := range payloads {
		_, err := event.Decode(payload)
		if err == nil {
			t.Fatal("Decode() error = nil, want invalid-event error")
		}
		if strings.Contains(err.Error(), payloadMarker) {
			t.Errorf("Decode() error %q contains telemetry payload data", err)
		}
	}
}

func FuzzDecode(f *testing.F) {
	f.Add([]byte(`{"id":"1","timestamp":"2026-09-19T18:30:45Z","actor":"a","action":"b","target":"c","source":"d"}`))
	f.Add([]byte(`{"id":`))
	f.Add([]byte(`[]`))

	f.Fuzz(func(t *testing.T, payload []byte) {
		record, err := event.Decode(payload)
		if err != nil {
			return
		}

		if got := record.Raw(); !bytes.Equal(got, payload) {
			t.Fatalf("Decode() Raw = %q, want exact payload %q", got, payload)
		}
	})
}
