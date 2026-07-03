package engine

import (
	"fmt"
	"sync"
	"time"
)

type IntegrationStore struct {
	mu           sync.RWMutex
	integrations map[string]*Integration
	counter      int
}

func NewIntegrationStore() *IntegrationStore {
	return &IntegrationStore{
		integrations: make(map[string]*Integration),
	}
}

func (s *IntegrationStore) Register(userID, name, provider, apiKey, webhookURL string) (*Integration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	id := fmt.Sprintf("int-%d", s.counter)
	now := time.Now().UTC().Format(time.RFC3339)
	integ := &Integration{
		ID: id, UserID: userID, Name: name, Provider: provider,
		APIKey: apiKey, WebhookURL: webhookURL, Enabled: true,
		CreatedAt: now,
	}
	s.integrations[id] = integ
	return integ, nil
}

func (s *IntegrationStore) ListByUser(userID string) ([]*Integration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Integration
	for _, i := range s.integrations {
		if i.UserID == userID {
			result = append(result, i)
		}
	}
	return result, nil
}

func (s *IntegrationStore) GetByID(id string) (*Integration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.integrations[id]
	if !ok {
		return nil, fmt.Errorf("integration not found: %s", id)
	}
	return i, nil
}

func (s *IntegrationStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.integrations[id]; !ok {
		return fmt.Errorf("integration not found: %s", id)
	}
	delete(s.integrations, id)
	return nil
}

func (s *IntegrationStore) Toggle(id string) (*Integration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, ok := s.integrations[id]
	if !ok {
		return nil, fmt.Errorf("integration not found: %s", id)
	}
	i.Enabled = !i.Enabled
	return i, nil
}
