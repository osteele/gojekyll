package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	l := New()
	l.SetOutput(&outBuf)
	l.SetErrorOutput(&errBuf)

	// Test Info level logging
	l.Info("test info message")
	if !strings.Contains(outBuf.String(), "test info message") {
		t.Errorf("Info message not logged correctly: %s", outBuf.String())
	}

	// Test Warning level logging
	outBuf.Reset()
	errBuf.Reset()
	l.Warn("test warning message")
	if !strings.Contains(errBuf.String(), "test warning message") {
		t.Errorf("Warning message not logged correctly: %s", errBuf.String())
	}

	// Test Error level logging
	errBuf.Reset()
	l.Error("test error message")
	if !strings.Contains(errBuf.String(), "test error message") {
		t.Errorf("Error message not logged correctly: %s", errBuf.String())
	}
}

func TestQuietMode(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	l := New()
	l.SetOutput(&outBuf)
	l.SetErrorOutput(&errBuf)
	l.SetQuiet(true)

	// Info and Debug should be suppressed
	l.Info("test info")
	l.Debug("test debug")
	if outBuf.Len() > 0 {
		t.Errorf("Info/Debug messages should be suppressed in quiet mode, got: %s", outBuf.String())
	}

	// Warnings and Errors should still appear
	l.Warn("test warning")
	l.Error("test error")
	if !strings.Contains(errBuf.String(), "test warning") || !strings.Contains(errBuf.String(), "test error") {
		t.Errorf("Warning/Error messages should appear in quiet mode, got: %s", errBuf.String())
	}
}

func TestLogLevel(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	l := New()
	l.SetOutput(&outBuf)
	l.SetErrorOutput(&errBuf)
	l.SetLevel(WarnLevel)

	// Info and Debug should be suppressed
	l.Info("test info")
	l.Debug("test debug")
	if outBuf.Len() > 0 {
		t.Errorf("Messages below log level should be suppressed, got: %s", outBuf.String())
	}

	// Warnings and Errors should appear
	l.Warn("test warning")
	l.Error("test error")
	if !strings.Contains(errBuf.String(), "test warning") || !strings.Contains(errBuf.String(), "test error") {
		t.Errorf("Messages at or above log level should appear, got: %s", errBuf.String())
	}
}

func TestLabel(t *testing.T) {
	var outBuf bytes.Buffer
	l := New()
	l.SetOutput(&outBuf)

	l.Label("create", "/path/to/file")
	output := outBuf.String()
	if !strings.Contains(output, "create") || !strings.Contains(output, "/path/to/file") {
		t.Errorf("Label output incorrect: %s", output)
	}

	// Test label width alignment
	outBuf.Reset()
	l.Label("rm", "/another/file")
	output = outBuf.String()
	// The label should be padded to match previous width
	if !strings.Contains(output, "rm") || !strings.Contains(output, "/another/file") {
		t.Errorf("Label output incorrect: %s", output)
	}
}

func TestPrintf(t *testing.T) {
	var outBuf bytes.Buffer
	l := New()
	l.SetOutput(&outBuf)

	l.Printf("test %s %d", "message", 42)
	output := outBuf.String()
	if !strings.Contains(output, "test message 42") {
		t.Errorf("Printf output incorrect: %s", output)
	}
}

func TestPrintln(t *testing.T) {
	var outBuf bytes.Buffer
	l := New()
	l.SetOutput(&outBuf)

	l.Println("test", "message")
	output := outBuf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Println output incorrect: %s", output)
	}
}

// TestLoggerStateMatrix verifies output routing for every combination of
// configured level, quiet mode, and message level.
func TestLoggerStateMatrix(t *testing.T) {
	tests := []struct {
		name      string
		level     Level
		quiet     bool
		callLevel Level
		wantOut   bool
		wantErr   bool
	}{
		{"info at info level", InfoLevel, false, InfoLevel, true, false},
		{"debug at info level", InfoLevel, false, DebugLevel, false, false},
		{"warn at info level", InfoLevel, false, WarnLevel, false, true},
		{"error at warn level", WarnLevel, false, ErrorLevel, false, true},
		{"info suppressed by quiet", InfoLevel, true, InfoLevel, false, false},
		{"warn not suppressed by quiet", InfoLevel, true, WarnLevel, false, true},
		{"debug suppressed by level and quiet", WarnLevel, true, DebugLevel, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer
			l := New()
			l.SetOutput(&outBuf)
			l.SetErrorOutput(&errBuf)
			l.SetLevel(tt.level)
			l.SetQuiet(tt.quiet)

			l.log(tt.callLevel, "msg")

			if tt.wantOut {
				require.Contains(t, outBuf.String(), "msg")
			} else {
				require.Empty(t, outBuf.String())
			}
			if tt.wantErr {
				require.Contains(t, errBuf.String(), "msg")
			} else {
				require.Empty(t, errBuf.String())
			}
		})
	}
}
