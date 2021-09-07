package domain

import "strings"

type Currency int

const (
	RateUndefined Currency = iota
	USD
	EUR
	PLN
	CAD
	UAH
)

func (c Currency) ToLabel() string {
	return [...]string{"", "USD", "EUR", "PLN", "CAD", "UAH"}[c]
}

func CurrencyToDomain(possibleCurrency string) Currency {
	currencyToUpper := strings.ToUpper(possibleCurrency)
	switch currencyToUpper {
	case "USD":
		return USD
	case "EUR":
		return EUR
	case "PLN":
		return PLN
	case "CAD":
		return CAD
	case "UAH":
		return UAH
	default:
		return RateUndefined
	}
}