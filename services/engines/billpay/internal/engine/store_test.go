package engine

import "testing"

func TestAddBiller(t *testing.T) {
	store := NewBillPayStore()
	b, err := store.AddBiller("user-1", "Netflix", "subscription", "12345", 1500, 15, true)
	if err != nil {
		t.Fatalf("add biller: %v", err)
	}
	if b.ID == "" { t.Error("expected non-empty id") }
	if b.Name != "Netflix" { t.Errorf("expected Netflix, got %s", b.Name) }
	if b.Amount != 1500 { t.Errorf("expected 1500, got %d", b.Amount) }
	if b.DueDay != 15 { t.Errorf("expected 15, got %d", b.DueDay) }
	if !b.AutoPay { t.Error("expected auto_pay true") }
}

func TestListBillers(t *testing.T) {
	store := NewBillPayStore()
	store.AddBiller("user-1", "Netflix", "", "", 1500, 15, false)
	store.AddBiller("user-1", "Electricity", "", "", 3000, 10, false)
	store.AddBiller("user-2", "Rent", "", "", 25000, 1, false)

	bills, _ := store.ListBillers("user-1")
	if len(bills) != 2 { t.Errorf("expected 2, got %d", len(bills)) }

	bills2, _ := store.ListBillers("user-2")
	if len(bills2) != 1 { t.Errorf("expected 1, got %d", len(bills2)) }
}

func TestPayBill(t *testing.T) {
	store := NewBillPayStore()
	b, _ := store.AddBiller("user-1", "Netflix", "", "", 1500, 15, false)

	p, err := store.PayBill(b.ID, "user-1")
	if err != nil { t.Fatalf("pay: %v", err) }
	if p.Status != "completed" { t.Errorf("expected completed, got %s", p.Status) }
	if p.Amount != 1500 { t.Errorf("expected 1500, got %d", p.Amount) }

	// Verify biller amount reset
	b2, _ := store.GetBiller(b.ID)
	if b2.Amount != 0 { t.Errorf("expected 0 after payment, got %d", b2.Amount) }
}

func TestDeleteBiller(t *testing.T) {
	store := NewBillPayStore()
	b, _ := store.AddBiller("user-1", "Test", "", "", 100, 1, false)
	if err := store.DeleteBiller(b.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := store.GetBiller(b.ID); err == nil {
		t.Error("expected error after delete")
	}
}

func TestPaymentHistory(t *testing.T) {
	store := NewBillPayStore()
	b, _ := store.AddBiller("user-1", "Netflix", "", "", 1500, 15, false)
	store.PayBill(b.ID, "user-1")

	payments, _ := store.ListPayments("user-1")
	if len(payments) != 1 { t.Errorf("expected 1 payment, got %d", len(payments)) }
}
