package testing

import (
	"github.com/horizon/core/packages/logging"
)

type MockLogger struct {
	entries []logEntry
}

type logEntry struct {
	Level  string
	Msg    string
	Fields []logging.Field
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (l *MockLogger) Debug(msg string, fields ...logging.Field) {
	l.entries = append(l.entries, logEntry{Level: "debug", Msg: msg, Fields: fields})
}

func (l *MockLogger) Info(msg string, fields ...logging.Field) {
	l.entries = append(l.entries, logEntry{Level: "info", Msg: msg, Fields: fields})
}

func (l *MockLogger) Warn(msg string, fields ...logging.Field) {
	l.entries = append(l.entries, logEntry{Level: "warn", Msg: msg, Fields: fields})
}

func (l *MockLogger) Error(msg string, fields ...logging.Field) {
	l.entries = append(l.entries, logEntry{Level: "error", Msg: msg, Fields: fields})
}

func (l *MockLogger) With(fields ...logging.Field) logging.Logger {
	return l
}

func (l *MockLogger) Entries() []logEntry {
	return l.entries
}

func (l *MockLogger) Clear() {
	l.entries = nil
}
