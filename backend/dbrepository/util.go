package dbrepository

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func isDup(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	if strings.Contains(err.Error(), "duplicate key") {
		return true
	}
	return false
}

func isFK(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return true
	}

	msg := err.Error()
	return strings.Contains(msg, "foreign key") ||
		strings.Contains(msg, "violates foreign key constraint")
}
