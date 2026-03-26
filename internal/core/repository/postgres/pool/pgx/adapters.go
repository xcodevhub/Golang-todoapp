package core_pgx_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/repository/postgres/pool"
)

type pgxRows struct {
	pgx.Rows
}

func (p pgxRows) Close() {
	p.Rows.Close()
}

func (p pgxRows) Err() error {
	return p.Rows.Err()
}

func (p pgxRows) Next() bool {
	return p.Rows.Next()
}

func (p pgxRows) Scan(dest ...any) error {
	err := p.Rows.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}
	return nil
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}
	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyErrorCode = "23503"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgxViolatesForeignKeyErrorCode {
			return fmt.Errorf(
				"%v: %w",
				err,
				core_postgres_pool.ErrViolatesForeignKey,
			)
		}
	}

	return fmt.Errorf(
		"%v: %w",
		err,
		core_postgres_pool.ErrUnknown,
	)
}
