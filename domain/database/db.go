package database

import (
	"context"
	"errors"
	"log"
	"passbase-hometest/config"
	"passbase-hometest/domain"
	"passbase-hometest/domain/database/sqlite"
)

type Repository interface {
	CreateRate(ctx context.Context, rate domain.Rate) (int64, error)
	UpdateRate(ctx context.Context, rate domain.Rate) error
	GetCurrencyRate(ctx context.Context, currency string) (*domain.Rate, error)
	GetBaseCurrency(ctx context.Context) (*domain.Rate, error)

	RegisterProject(ctx context.Context, project domain.Project) (int64, error)
	FindProjectByEmail(ctx context.Context, email string) (*domain.Project, error)
	FindProjectByToken(ctx context.Context, token string) (*domain.Project, error)
}

var ErrRecordNotFound = errors.New("record not found")

type Databases int
const (
	SQLite Databases = iota+1
)

func Connect(conf config.Database, dbtype Databases) Repository {
	switch dbtype {
	case SQLite:
		return sqlite.New(conf)
	default:
		log.Fatalf("implementation missing for DB type: %d", dbtype)
		return nil
	}
}