package tasks_service

import (
	"context"
	"fmt"

	"github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/domain"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	id int,
) (*domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return &domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}
	return &task, nil

}
