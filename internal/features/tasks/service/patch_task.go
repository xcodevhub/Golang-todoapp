package tasks_service

import (
	"context"
	"fmt"

	"github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/domain"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	id int,
	patch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}
	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply patch: %w", err)
	}
	patchedTask, err := s.tasksRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}
	return patchedTask, nil
}
