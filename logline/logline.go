// Package logline buffers a single line of log from multiple calls
package logline

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
)

// LogLine buffers a single line of log.
//
// Each Print() call append a string to the buffer and Write()
// will write the full line to the log.
type LogLine struct {
	buf    *bytes.Buffer
	sep    string
	logger *log.Logger
	mu     sync.Mutex
}

type Info struct {
	// A string that will be inserted automatically between Print()
	// calls.
	// It's ". " by default.
	Separator string

	// Logger where the line will be written to.
	// Stderr by default (same as log package)
	Logger *log.Logger
}

// New creates a new log line buffer. Call Print() functions then
// Write() to actuall write the log.
//
// If no logger is configured with Info (see NewWithInfo) will log to
// stderr by default (same as log package)
func New() *LogLine {
	return NewWithInfo(Info{})
}

// NewWithInfo creates a new customized log line buffer.
func NewWithInfo(info Info) *LogLine {
	if info.Separator == "" {
		info.Separator = ". "
	}
	if info.Logger == nil {
		info.Logger = log.New(os.Stderr, "", log.LstdFlags) // same as 'std' on log package
	}
	return &LogLine{
		buf:    &bytes.Buffer{},
		sep:    info.Separator,
		logger: info.Logger,
	}
}

// Print appends strings to the line buffer.
// Arguments are handled in the manner of fmt.Print.
func (ll *LogLine) Printf(format string, a ...interface{}) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	ll.buf.WriteString(fmt.Sprintf(format, a...))
	ll.buf.WriteString(ll.sep)
}

// Print appends strings to the line buffer.
// Arguments are handled in the manner of fmt.Print.
func (ll *LogLine) Print(v ...interface{}) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	ll.buf.WriteString(fmt.Sprint(v...))
	ll.buf.WriteString(ll.sep)
}

// Write writes line to the log and clears the buffer.
func (ll *LogLine) Write() {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	ll.logger.Print(ll.buf.String())
	ll.buf.Reset()
}
