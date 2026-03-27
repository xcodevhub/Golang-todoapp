package statistics_postgres_repository

import core_postgres_pool "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/repository/postgres/pool"

type StatisticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatisticsRepository(
	pool core_postgres_pool.Pool,

) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
