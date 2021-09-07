package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"passbase-hometest/domain"
	"time"
)

type Rate struct {
	ID              int64      `db:"id"`
	Rate    float64   `db:"rate"`
	Currency string    `db:"currency"`
	Updated  time.Time `db:"updated"`
	Base     bool      `db:"base"`
}

func (s *SQLite) CreateRate(ctx context.Context, rate domain.Rate) (int64, error) {
	stmt, err := s.connection.PrepareContext(ctx, "INSERT INTO rates(currency, rate, base, updated) values(?,?,?,?)")
	if err != nil {
		logger.Errorf("can't prepare data to insert rate: err: %s", err)
		return -1, err
	}

	res, err := stmt.Exec(rate.Currency, rate.Value, rate.Base, rate.Updated)
	if err != nil {
		logger.Errorf("can't insert rate: err: %s", err)
		return -1, err
	}

	return res.LastInsertId()
}

func (s *SQLite) UpdateRate(ctx context.Context, rate domain.Rate) error {
	stmt, err := s.connection.PrepareContext(ctx, "update rates set rate=?, base=?, updated=? where currency=?")
	if err != nil {
		logger.Errorf("can't prepare data to update rate: err: %s", err)
		return err
	}

	_, err = stmt.Exec(rate.Value, rate.Base, rate.Updated, rate.Currency)
	if err != nil {
		logger.Errorf("can't update rate: err: %s", err)
		return err
	}

	return nil
}

func (s *SQLite) GetCurrencyRate(ctx context.Context, currency string) (*domain.Rate, error) {
	var rate Rate
	err := s.connection.
		QueryRowContext(ctx, "SELECT currency, rate, base, updated FROM rates where currency=?", currency).
		Scan(&rate.Currency, &rate.Rate, &rate.Base, &rate.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Errorf("can't execute query to get rate: err: %s", err)
		return nil, err
	}

	return &domain.Rate{
		Value:    rate.Rate,
		Domain:   domain.CurrencyToDomain(rate.Currency),
		Currency: rate.Currency,
		Updated:  rate.Updated,
		Base:     rate.Base,
	}, nil
}

func (s *SQLite) GetBaseCurrency(ctx context.Context) (*domain.Rate, error) {
	var rate Rate
	err := s.connection.
		QueryRowContext(ctx, "SELECT currency, rate, base, updated FROM rates where base=true").
		Scan(&rate.Currency, &rate.Rate, &rate.Base, &rate.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Errorf("can't execute query to get rate: err: %s", err)
		return nil, err
	}

	return &domain.Rate{
		Value:    rate.Rate,
		Domain:   domain.CurrencyToDomain(rate.Currency),
		Currency: rate.Currency,
		Updated:  rate.Updated,
		Base:     rate.Base,
	}, nil
}