package engine

import "testing"

func TestGetCurrencies(t *testing.T) {
	currencies := GetCurrencies()
	if len(currencies) == 0 {
		t.Fatal("expected non-empty currency list")
	}
	foundINR := false
	for _, c := range currencies {
		if c.Code == "INR" {
			foundINR = true
			if c.Rate != 1.0 {
				t.Errorf("expected INR rate 1.0, got %f", c.Rate)
			}
			break
		}
	}
	if !foundINR {
		t.Error("INR not found in currencies")
	}
}

func TestGetRate(t *testing.T) {
	rate, err := GetRate("USD")
	if err != nil {
		t.Fatalf("get USD rate: %v", err)
	}
	if rate <= 0 {
		t.Errorf("expected positive rate, got %f", rate)
	}

	_, err = GetRate("XYZ")
	if err == nil {
		t.Error("expected error for unsupported currency")
	}
}

func TestConvert(t *testing.T) {
	result, err := Convert(ConvertRequest{From: "USD", To: "INR", Amount: 100})
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if result.From != "USD" { t.Errorf("expected USD from, got %s", result.From) }
	if result.To != "INR" { t.Errorf("expected INR to, got %s", result.To) }
	if result.Amount != 100 { t.Errorf("expected 100 amount, got %f", result.Amount) }
	if result.Result <= 0 { t.Errorf("expected positive result, got %f", result.Result) }
	if result.Rate <= 0 { t.Errorf("expected positive rate, got %f", result.Rate) }
}

func TestConvert_SameCurrency(t *testing.T) {
	result, err := Convert(ConvertRequest{From: "INR", To: "INR", Amount: 1000})
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if result.Result != 1000 {
		t.Errorf("expected 1000 for same currency, got %f", result.Result)
	}
	if result.Rate != 1.0 {
		t.Errorf("expected rate 1.0 for same currency, got %f", result.Rate)
	}
}

func TestConvertToINR(t *testing.T) {
	amount, err := ConvertToINR(100, "USD")
	if err != nil {
		t.Fatalf("convert to INR: %v", err)
	}
	if amount <= 0 {
		t.Errorf("expected positive amount, got %f", amount)
	}
}

func TestConvertFromINR(t *testing.T) {
	amount, err := ConvertFromINR(8300, "USD")
	if err != nil {
		t.Fatalf("convert from INR: %v", err)
	}
	if amount <= 0 {
		t.Errorf("expected positive amount, got %f", amount)
	}
}

func TestUnsupportedCurrency(t *testing.T) {
	_, err := Convert(ConvertRequest{From: "XYZ", To: "INR", Amount: 100})
	if err == nil {
		t.Error("expected error for unsupported currency")
	}
}
