package domain

import (
	"time"
)

type Rate struct {
	Value    float64
	Domain   Currency
	Currency string
	Updated  time.Time
	Base     bool
}

// Convert serves the idea of money conversion between currencies,
// should be called ONLY as a method of *base rate* to convert prices correctly
func (source *Rate) Convert(destination Rate, amount int) float64 {
	convertedToBaseAmount := float64(amount) / source.Value
	convertedDestination := convertedToBaseAmount * destination.Value

	return convertedDestination
}
