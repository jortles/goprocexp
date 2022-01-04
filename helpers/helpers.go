package helpers

import (
	"sync"
	"time"
)

type Logger struct {
	mu      sync.Mutex
	entries ProcessEntries
}

type Notification struct {
	ID        string
	NotifType int
}

type ProcessEntries map[string]*Entry

type Entry struct {
	ID              string    `json:"_id"`
	StartedDateTime time.Time `json:"startedDateTime"`
	Time            int64     `json:"time"`

	Details *ProcessDetails `json:"details"`
}

type ProcessDetails struct {
	Status        string `json:"status"`
	MemoryPercent int    `json:"memoryPercent"`
	CommandLine   string `json:"commandLine"`
	ThreadCount   int    `json:"threadCount"`

	Raw []byte
}

func (l *Logger) GetEntries() map[string]*Entry {
	return l.entries
}

func (l *Logger) AddEntry(e Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e.ID != "" {
		l.entries[e.ID] = &e
	}
}

func (l *Logger) GetEntry(id string) *Entry {
	var e *Entry

	e = l.entries[id]

	return e
}

func (l *Logger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = make(map[string]*Entry)
}
