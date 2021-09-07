package providers

import (
	"errors"
	"passbase-hometest/config"
	"passbase-hometest/domain"
	"passbase-hometest/domain/providers/fixer"
)

type Provider int

const (
	ProviderUndefined Provider = iota
	Fixer
)

// Providerer is not the best name, I know :)
type Providerer interface {
	// Latest gets the latest rate results from the provider
	// returns:
	// [1] - list of all possible rates
	// [2] - error if occurs
	Latest() ([]domain.Rate, error)
}

func New(provider Provider, conf *config.Config) (Providerer, error) {
	switch provider {
	case Fixer:
		return fixer.New(&conf.Fixer), nil
	default:
		return nil, errors.New("rates provider is not found")
	}
}
