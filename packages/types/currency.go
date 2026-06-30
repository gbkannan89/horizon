package types

import "fmt"

type Currency struct {
	Code           string
	Symbol         string
	DecimalPlaces  int
	ThousandSep    string
}

func INR() Currency {
	return Currency{Code: "INR", Symbol: "₹", DecimalPlaces: 2, ThousandSep: ","}
}

func USD() Currency {
	return Currency{Code: "USD", Symbol: "$", DecimalPlaces: 2, ThousandSep: ","}
}

func EUR() Currency {
	return Currency{Code: "EUR", Symbol: "€", DecimalPlaces: 2, ThousandSep: "."}
}

func GBP() Currency {
	return Currency{Code: "GBP", Symbol: "£", DecimalPlaces: 2, ThousandSep: ","}
}

var currencies = map[string]Currency{
	"INR": INR(),
	"USD": USD(),
	"EUR": EUR(),
	"GBP": GBP(),
}

func GetCurrency(code string) (Currency, error) {
	c, ok := currencies[code]
	if !ok {
		return Currency{}, fmt.Errorf("unsupported currency: %s", code)
	}
	return c, nil
}

func SupportedCurrencies() []string {
	keys := make([]string, 0, len(currencies))
	for k := range currencies {
		keys = append(keys, k)
	}
	return keys
}
