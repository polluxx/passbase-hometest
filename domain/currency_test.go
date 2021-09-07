package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrencyToDomain(t *testing.T) {
	tests := []struct {
		name        string
		currencyStr string
		currency    Currency
	}{
		{
			name:        "should return undefined on incorrect currency",
			currencyStr: "bleh",
			currency:    RateUndefined,
		},
		{
			name:        "should pass on correct currency: USD",
			currencyStr: "USD",
			currency:    USD,
		},
		{
			name:        "should pass on correct currency lowercase: eur",
			currencyStr: "eur",
			currency:    EUR,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CurrencyToDomain(tc.currencyStr)
			assert.Equal(t, result, tc.currency, "should be equal")
		})
	}
}
