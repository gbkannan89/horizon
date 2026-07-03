package engine

import "fmt"

var currencies = []CurrencyInfo{
	{Code: "INR", Name: "Indian Rupee", Symbol: "₹", Rate: 1.0},
	{Code: "USD", Name: "US Dollar", Symbol: "$", Rate: 83.0},
	{Code: "EUR", Name: "Euro", Symbol: "€", Rate: 90.5},
	{Code: "GBP", Name: "British Pound", Symbol: "£", Rate: 105.0},
	{Code: "CAD", Name: "Canadian Dollar", Symbol: "C$", Rate: 61.0},
	{Code: "AUD", Name: "Australian Dollar", Symbol: "A$", Rate: 55.0},
	{Code: "SGD", Name: "Singapore Dollar", Symbol: "S$", Rate: 62.0},
	{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", Rate: 0.56},
	{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", Rate: 11.5},
	{Code: "CHF", Name: "Swiss Franc", Symbol: "Fr", Rate: 94.0},
	{Code: "AED", Name: "UAE Dirham", Symbol: "د.إ", Rate: 22.6},
	{Code: "SAR", Name: "Saudi Riyal", Symbol: "﷼", Rate: 22.1},
}

func GetCurrencies() []CurrencyInfo {
	return currencies
}

func GetRate(code string) (float64, error) {
	for _, c := range currencies {
		if c.Code == code {
			return c.Rate, nil
		}
	}
	return 0, fmt.Errorf("unsupported currency: %s", code)
}

func Convert(req ConvertRequest) (*ConvertResult, error) {
	fromRate, err := GetRate(req.From)
	if err != nil {
		return nil, fmt.Errorf("from currency: %w", err)
	}
	toRate, err := GetRate(req.To)
	if err != nil {
		return nil, fmt.Errorf("to currency: %w", err)
	}

	inINR := req.Amount * fromRate
	result := inINR / toRate

	return &ConvertResult{
		From:   req.From, To: req.To,
		Amount: req.Amount, Result: result,
		Rate: fromRate / toRate,
	}, nil
}

func ConvertToINR(amount float64, from string) (float64, error) {
	rate, err := GetRate(from)
	if err != nil {
		return 0, err
	}
	return amount * rate, nil
}

func ConvertFromINR(amount float64, to string) (float64, error) {
	rate, err := GetRate(to)
	if err != nil {
		return 0, err
	}
	return amount / rate, nil
}
