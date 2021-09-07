package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"passbase-hometest/domain"
)

type Project struct {
	ID            int64  `db:"id"`
	Name          string `db:"name"`
	CustomerEmail string `db:"email"`
	Token         string `db:"token"`
}

func (s *SQLite) RegisterProject(ctx context.Context, project domain.Project) (int64, error) {
	stmt, err := s.connection.PrepareContext(ctx, "INSERT INTO projects(name, email, token) values(?,?,?)")
	if err != nil {
		logger.Errorf("can't prepare data to insert project: err: %s", err)
		return -1, err
	}

	res, err := stmt.Exec(project.Name, project.CustomerEmail, project.Token)
	if err != nil {
		logger.Errorf("can't insert project: err: %s", err)
		return -1, err
	}

	return res.LastInsertId()
}

func (s *SQLite) FindProjectByEmail(ctx context.Context, email string) (*domain.Project, error) {
	var project Project
	err := s.connection.
		QueryRowContext(ctx, "SELECT name, email, token FROM projects where email=?", email).
		Scan(&project.Name, &project.CustomerEmail, &project.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Errorf("can't execute query to get project: err: %s", err)
		return nil, err
	}

	return &domain.Project{
		Name:          project.Name,
		CustomerEmail: project.CustomerEmail,
		Token:         project.Token,
	}, nil
}

func (s *SQLite) FindProjectByToken(ctx context.Context, token string) (*domain.Project, error) {
	var project Project
	err := s.connection.
		QueryRowContext(ctx, "SELECT name, email, token FROM projects where token=?", token).
		Scan(&project.Name, &project.CustomerEmail, &project.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Errorf("can't execute query to get project by token: err: %s", err)
		return nil, err
	}

	return &domain.Project{
		Name:          project.Name,
		CustomerEmail: project.CustomerEmail,
		Token:         project.Token,
	}, nil
}
