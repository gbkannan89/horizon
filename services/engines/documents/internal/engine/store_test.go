package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndGet(t *testing.T) {
	dir := t.TempDir()
	store := NewDocumentStore(dir)

	doc, err := store.Save("user-1", "test.pdf", "Test doc", "bank-statement", "application/pdf", []byte("test content"), []string{"important"})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if doc.ID == "" {
		t.Error("expected non-empty id")
	}
	if doc.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", doc.UserID)
	}
	if doc.Name != "test.pdf" {
		t.Errorf("expected test.pdf, got %s", doc.Name)
	}
	if doc.FileSize != 12 {
		t.Errorf("expected 12, got %d", doc.FileSize)
	}
	if doc.MimeType != "application/pdf" {
		t.Errorf("expected application/pdf, got %s", doc.MimeType)
	}
	if len(doc.Tags) != 1 || doc.Tags[0] != "important" {
		t.Errorf("expected [important], got %v", doc.Tags)
	}

	// Verify file was saved
	filePath := filepath.Join(dir, doc.ID+"-"+doc.Name)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("saved file not found on disk")
	}

	// Get by ID
	got, err := store.GetByID(doc.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "test.pdf" {
		t.Errorf("expected test.pdf, got %s", got.Name)
	}
}

func TestListByUser(t *testing.T) {
	dir := t.TempDir()
	store := NewDocumentStore(dir)

	store.Save("user-1", "a.pdf", "", "", "", []byte("a"), nil)
	store.Save("user-1", "b.pdf", "", "", "", []byte("b"), nil)
	store.Save("user-2", "c.pdf", "", "", "", []byte("c"), nil)

	docs, err := store.ListByUser("user-1")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("expected 2 docs for user-1, got %d", len(docs))
	}

	docs2, err := store.ListByUser("user-2")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(docs2) != 1 {
		t.Errorf("expected 1 doc for user-2, got %d", len(docs2))
	}

	docs3, err := store.ListByUser("user-3")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(docs3) != 0 {
		t.Errorf("expected 0 docs for user-3, got %d", len(docs3))
	}
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	store := NewDocumentStore(dir)
	content := []byte("hello world")

	doc, err := store.Save("user-1", "hello.txt", "", "", "text/plain", content, nil)
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	data, err := store.ReadFile(doc.ID)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(data))
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	store := NewDocumentStore(dir)

	doc, _ := store.Save("user-1", "del.txt", "", "", "", []byte("data"), nil)

	if err := store.Delete(doc.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := store.GetByID(doc.ID)
	if err == nil {
		t.Error("expected error after delete")
	}

	filePath := filepath.Join(dir, doc.ID+"-"+doc.Name)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("file should be deleted from disk")
	}
}

func TestDeleteNotFound(t *testing.T) {
	dir := t.TempDir()
	store := NewDocumentStore(dir)
	if err := store.Delete("nonexistent"); err == nil {
		t.Error("expected error for non-existent id")
	}
}

func TestGetByIDNotFound(t *testing.T) {
	dir := t.TempDir()
	store := NewDocumentStore(dir)
	_, err := store.GetByID("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent id")
	}
}
