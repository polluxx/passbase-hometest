package fixer

import (
	"encoding/json"
	"fmt"
	"passbase-hometest/config"
	"passbase-hometest/domain"
	"passbase-hometest/domain/http"
	"time"
)

type Fixer struct {
	latestRatesEndpoint string
	apikey              string
	basehost            string
}

type Response struct {
	Success   bool   `json:"success"`
	Base      string `json:"base"`
	Timestamp int64  `json:"timestamp"`
	Rates     map[string]float64
}

func New(conf *config.Fixer) *Fixer {
	return &Fixer{
		latestRatesEndpoint: "api/latest",
		apikey:              conf.APIKey,
		basehost:            conf.BaseHost,
	}
}

func (f *Fixer) Latest() ([]domain.Rate, error) {
	url := fmt.Sprintf("%s/%s", f.basehost, f.latestRatesEndpoint)
	httpClient := &http.Request{Debug: true}

	body, status, err := httpClient.Dial().Get(url, []http.QueryParam{{Key: "access_key", Value: f.apikey}})
	if err != nil {
		return nil, fmt.Errorf("can't call rates endpoint: %w", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("rates endpoint status is: %d; response: %s", status, string(body))
	}

	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("can't unmarshal rates response: %w; payload: %s", err, string(body))
	}

	rates := make([]domain.Rate, 0)
	updateTime := time.Unix(response.Timestamp, 0)
	for currency, value := range response.Rates {
		domainRate := domain.CurrencyToDomain(currency)
		// filter out all unsupported currencies
		if domainRate == domain.RateUndefined {
			continue
		}

		rates = append(rates, domain.Rate{
			Value:    value,
			Domain:   domainRate,
			Currency: currency,
			Base:     response.Base == currency,
			Updated:  updateTime,
		})
	}

	return rates, nil
}
