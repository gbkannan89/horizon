package types

import (
	"fmt"
	"strings"
)

type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, fmt.Errorf("amount must be non-negative: %d", amount)
	}
	if currency == "" {
		return Money{}, fmt.Errorf("currency is required")
	}
	currency = strings.ToUpper(currency)
	return Money{amount: amount, currency: currency}, nil
}

func MustMoney(amount int64, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }
func (m Money) IsZero() bool     { return m.amount == 0 }
func (m Money) IsPositive() bool { return m.amount > 0 }

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.currency, other.currency)
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.currency, other.currency)
	}
	if m.amount < other.amount {
		return Money{}, fmt.Errorf("insufficient funds: %d < %d", m.amount, other.amount)
	}
	return Money{amount: m.amount - other.amount, currency: m.currency}, nil
}

func (m Money) Multiply(factor int64) Money {
	return Money{amount: m.amount * factor, currency: m.currency}
}

func (m Money) Compare(other Money) int {
	if m.currency != other.currency {
		if m.currency < other.currency {
			return -1
		}
		return 1
	}
	switch {
	case m.amount < other.amount:
		return -1
	case m.amount > other.amount:
		return 1
	default:
		return 0
	}
}

func (m Money) ToFloat() float64 {
	return float64(m.amount) / 100
}

func (m Money) String() string {
	return fmt.Sprintf("%s %.2f", m.currency, m.ToFloat())
}
