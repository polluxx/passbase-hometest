package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRate_Convert(t *testing.T) {
	tests := []struct {
		name      string
		source    Rate
		dest      Rate
		amount    int
		converted float64
	}{
		{
			name: "should convert usd to uah",
			source: Rate{
				Value:    1.2,
				Domain:   USD,
				Currency: "USD",
			},
			dest: Rate{
				Value:    30.0,
				Domain:   UAH,
				Currency: "UAH",
			},
			amount:    120,
			converted: 3000,
		},
		{
			name: "should convert eur to cad",
			source: Rate{
				Value:    1.0,
				Domain:   EUR,
				Currency: "EUR",
			},
			dest: Rate{
				Value:    1.23,
				Domain:   CAD,
				Currency: "CAD",
			},
			amount:    150,
			converted: 184.5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.source.Convert(tc.dest, tc.amount)
			assert.Equal(t, tc.converted, result, "should be equal")
		})
	}
}
