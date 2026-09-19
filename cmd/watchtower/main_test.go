package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

var errWrite = errors.New("write failed")

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

func TestRunVersion(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	got := run([]string{"--version"}, &stdout, &stderr)

	if got != exitSuccess {
		t.Fatalf("run() exit code = %d, want %d", got, exitSuccess)
	}
	if want := "watchtower dev\n"; stdout.String() != want {
		t.Errorf("run() stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Errorf("run() stderr = %q, want empty output", stderr.String())
	}
}

func TestRunRejectsUnexpectedArguments(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	got := run([]string{"unexpected"}, &stdout, &stderr)

	if got != exitUsage {
		t.Fatalf("run() exit code = %d, want %d", got, exitUsage)
	}
	if stdout.Len() != 0 {
		t.Errorf("run() stdout = %q, want empty output", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unexpected arguments") {
		t.Errorf("run() stderr = %q, want an unexpected-arguments error", stderr.String())
	}
}

func TestRunHelpSucceeds(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	got := run([]string{"--help"}, &stdout, &stderr)

	if got != exitSuccess {
		t.Fatalf("run() exit code = %d, want %d", got, exitSuccess)
	}
	if stdout.Len() != 0 {
		t.Errorf("run() stdout = %q, want empty output", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Usage: watchtower") {
		t.Errorf("run() stderr = %q, want usage output", stderr.String())
	}
}

func TestRunFailsWhenVersionCannotBeWritten(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer

	got := run([]string{"--version"}, errorWriter{}, &stderr)

	if got != exitFailure {
		t.Fatalf("run() exit code = %d, want %d", got, exitFailure)
	}
}
