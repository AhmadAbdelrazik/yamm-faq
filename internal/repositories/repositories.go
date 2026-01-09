package repositories

import (
	"database/sql"
	"errors"
	"log/slog"
)

var (
	ErrDuplicate = errors.New("duplicate entry")
	ErrNotFound  = errors.New("resource not found")
)

type Repositories struct {
	Users *UserRepository
}

// New creates new repository instance using dsn.
func New(dsn string) (*Repositories, error) {
	slog.Debug("Connecting to database")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}
	slog.Debug(
		"Connection to database established",
		slog.Group(
			"database",
			slog.String("database", "postgres"),
			slog.String("port", "5432"),
		),
	)

	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", slog.String("error", err.Error()))
		return nil, err
	}

	return &Repositories{
		Users: &UserRepository{db},
	}, nil
}
