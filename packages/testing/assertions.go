package testing

import (
	"testing"

	"github.com/horizon/core/packages/types"
)

func AssertMoneyEqual(t *testing.T, expected, actual types.Money) {
	t.Helper()
	if expected.Amount() != actual.Amount() || expected.Currency() != actual.Currency() {
		t.Errorf("Money mismatch: expected %v, got %v", expected, actual)
	}
}

func AssertMoneyZero(t *testing.T, m types.Money) {
	t.Helper()
	if !m.IsZero() {
		t.Errorf("expected zero Money, got %v", m)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
