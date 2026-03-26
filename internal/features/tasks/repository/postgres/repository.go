package tasks_postgres_repository

import core_postgres_pool "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/repository/postgres/pool"

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(
	pool core_postgres_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
