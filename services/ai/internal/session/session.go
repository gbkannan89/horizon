package session

import (
	"sync"
	"time"
)

// Conversation represents a conversation session.
type Conversation struct {
	SessionID string          `json:"session_id"`
	UserID    string          `json:"user_id"`
	Messages  []Message       `json:"messages"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

type Message struct {
	Role      string `json:"role"` // user, assistant, system
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Confidence string `json:"confidence,omitempty"`
}

// Manager manages conversation sessions in-memory.
type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*Conversation
	maxHist  int // max messages to keep
}

func NewManager(maxHistory int) *Manager {
	return &Manager{sessions: make(map[string]*Conversation), maxHist: maxHistory}
}

func (m *Manager) GetOrCreate(sessionID, userID string) *Conversation {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		return s
	}
	now := time.Now().UTC().Format(time.RFC3339)
	s := &Conversation{
		SessionID: sessionID,
		UserID:    userID,
		Messages:  []Message{},
		CreatedAt: now,
		UpdatedAt: now,
		Context:   make(map[string]interface{}),
	}
	m.sessions[sessionID] = s
	return s
}

func (m *Manager) AddMessage(sessionID string, msg Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[sessionID]
	if !ok { return }
	s.Messages = append(s.Messages, msg)
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if len(s.Messages) > m.maxHist {
		s.Messages = s.Messages[len(s.Messages)-m.maxHist:]
	}
}

func (m *Manager) GetHistory(sessionID string) []Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[sessionID]
	if !ok { return nil }
	cp := make([]Message, len(s.Messages))
	copy(cp, s.Messages)
	return cp
}

func (m *Manager) SetContext(sessionID string, ctx map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		s.Context = ctx
	}
}

func (m *Manager) Get(sessionID string) *Conversation {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[sessionID]
}
