package engine

import (
	"fmt"
	"sync"
	"time"
)

type BillPayStore struct {
	mu         sync.RWMutex
	billers    map[string]*Biller
	payments   map[string]*Payment
	counter    int
}

func NewBillPayStore() *BillPayStore {
	return &BillPayStore{
		billers:  make(map[string]*Biller),
		payments: make(map[string]*Payment),
	}
}

func (s *BillPayStore) AddBiller(userID, name, category, accountNumber string, amount int64, dueDay int, autoPay bool) (*Biller, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	id := fmt.Sprintf("biller-%d", s.counter)
	b := &Biller{
		ID: id, UserID: userID, Name: name, Category: category,
		AccountNumber: accountNumber, Amount: amount, DueDay: dueDay,
		AutoPay: autoPay, Active: true,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.billers[id] = b
	return b, nil
}

func (s *BillPayStore) ListBillers(userID string) ([]*Biller, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Biller
	for _, b := range s.billers {
		if b.UserID == userID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (s *BillPayStore) GetBiller(id string) (*Biller, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.billers[id]
	if !ok {
		return nil, fmt.Errorf("biller not found: %s", id)
	}
	return b, nil
}

func (s *BillPayStore) DeleteBiller(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.billers[id]; !ok {
		return fmt.Errorf("biller not found: %s", id)
	}
	delete(s.billers, id)
	return nil
}

func (s *BillPayStore) PayBill(billerID, userID string) (*Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.billers[billerID]
	if !ok {
		return nil, fmt.Errorf("biller not found: %s", billerID)
	}
	s.counter++
	id := fmt.Sprintf("payment-%d", s.counter)
	p := &Payment{
		ID: id, BillerID: billerID, UserID: userID,
		Amount: b.Amount, Status: "completed",
		PaidAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.payments[id] = p
	s.counter++
	b.Amount = 0
	return p, nil
}

func (s *BillPayStore) ListPayments(userID string) ([]*Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Payment
	for _, p := range s.payments {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}
