package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DocumentStore struct {
	mu         sync.RWMutex
	documents  map[string]*Document
	uploadDir  string
	counter    int
}

func NewDocumentStore(uploadDir string) *DocumentStore {
	os.MkdirAll(uploadDir, 0755)
	return &DocumentStore{
		documents: make(map[string]*Document),
		uploadDir: uploadDir,
	}
}

func (s *DocumentStore) Save(userID, name, desc, category, mimeType string, data []byte, tags []string) (*Document, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	id := fmt.Sprintf("doc-%d", s.counter)
	fileName := fmt.Sprintf("%s-%s", id, name)
	filePath := filepath.Join(s.uploadDir, fileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}

	if desc == "" { desc = name }
	if category == "" { category = "other" }
	if tags == nil { tags = []string{} }

	now := time.Now().UTC().Format(time.RFC3339)
	doc := &Document{
		ID: id, UserID: userID, Name: name,
		Description: desc, Category: category,
		FileSize: int64(len(data)), MimeType: mimeType,
		Tags: tags, CreatedAt: now, UpdatedAt: now,
	}
	s.documents[id] = doc
	return doc, nil
}

func (s *DocumentStore) GetByID(id string) (*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.documents[id]
	if !ok {
		return nil, fmt.Errorf("document not found: %s", id)
	}
	return doc, nil
}

func (s *DocumentStore) ListByUser(userID string) ([]*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Document
	for _, d := range s.documents {
		if d.UserID == userID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (s *DocumentStore) ReadFile(id string) ([]byte, error) {
	doc, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filepath.Join(s.uploadDir, fmt.Sprintf("%s-%s", doc.ID, doc.Name)))
}

func (s *DocumentStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	doc, ok := s.documents[id]
	if !ok {
		return fmt.Errorf("document not found: %s", id)
	}
	os.Remove(filepath.Join(s.uploadDir, fmt.Sprintf("%s-%s", doc.ID, doc.Name)))
	delete(s.documents, id)
	return nil
}
