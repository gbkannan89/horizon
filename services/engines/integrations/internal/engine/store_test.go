package engine

import "testing"

func TestRegister(t *testing.T) {
	store := NewIntegrationStore()
	integ, err := store.Register("user-1", "My Bank", "plaid", "key-123", "https://hook.example.com")
	if err != nil { t.Fatalf("register: %v", err) }
	if integ.ID == "" { t.Error("expected non-empty id") }
	if integ.Name != "My Bank" { t.Errorf("expected My Bank, got %s", integ.Name) }
	if integ.Provider != "plaid" { t.Errorf("expected plaid, got %s", integ.Provider) }
	if !integ.Enabled { t.Error("expected enabled by default") }
}

func TestListIntegrations(t *testing.T) {
	store := NewIntegrationStore()
	store.Register("user-1", "Bank A", "plaid", "", "")
	store.Register("user-1", "Bank B", "finicity", "", "")
	store.Register("user-2", "Bank C", "yodlee", "", "")

	list, _ := store.ListByUser("user-1")
	if len(list) != 2 { t.Errorf("expected 2, got %d", len(list)) }
}

func TestToggle(t *testing.T) {
	store := NewIntegrationStore()
	integ, _ := store.Register("user-1", "Test", "test", "", "")

	toggled, err := store.Toggle(integ.ID)
	if err != nil { t.Fatalf("toggle: %v", err) }
	if toggled.Enabled { t.Error("expected disabled after toggle") }

	toggled2, _ := store.Toggle(integ.ID)
	if !toggled2.Enabled { t.Error("expected enabled after second toggle") }
}

func TestDelete(t *testing.T) {
	store := NewIntegrationStore()
	integ, _ := store.Register("user-1", "Test", "test", "", "")
	if err := store.Delete(integ.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := store.GetByID(integ.ID); err == nil {
		t.Error("expected error after delete")
	}
}
