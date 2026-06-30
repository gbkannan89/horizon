package types

import (
	"fmt"
	"math"
)

type Percentage struct {
	rate float64
}

func NewPercentage(rate float64) (Percentage, error) {
	if rate < 0 || rate > 100 {
		return Percentage{}, fmt.Errorf("percentage must be between 0 and 100: %f", rate)
	}
	return Percentage{rate: rate}, nil
}

func MustPercentage(rate float64) Percentage {
	p, err := NewPercentage(rate)
	if err != nil {
		panic(err)
	}
	return p
}

func (p Percentage) Rate() float64 { return p.rate }
func (p Percentage) Decimal() float64 { return p.rate / 100.0 }

func (p Percentage) Of(m Money) Money {
	amt := int64(math.Round(float64(m.Amount()) * p.Decimal()))
	return Money{amount: amt, currency: m.Currency()}
}

func (p Percentage) Add(other Percentage) Percentage {
	return Percentage{rate: p.rate + other.rate}
}

func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", p.rate)
}
