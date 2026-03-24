package core_pgx_pool

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/repository/postgres/pool"
)

type pgxRows struct {
	pgx.Rows
}

func (p pgxRows) Close() {
	panic("unimplemented")
}
func (p pgxRows) Err() {
	panic("unimplemented")
}
func (p pgxRows) Next() bool {
	panic("unimplemented")
}
func (p pgxRows) Scan(dest ...any) error {
	panic("unimplemented")
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_postgres_pool.ErrNoRows
		}
		return err
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}
