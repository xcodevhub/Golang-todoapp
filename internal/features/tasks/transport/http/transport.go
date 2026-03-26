package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/domain"
	core_http_server "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		userID *int,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		id int,
	) (*domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id int,
	) error

	PatchTask(
		ctx context.Context,
		id int,
		patch domain.TaskPatch,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(
	tasksService TasksService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: http.HandlerFunc(h.createTask),
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: http.HandlerFunc(h.GetTasks),
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: http.HandlerFunc(h.GetTask),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: http.HandlerFunc(h.DeleteTask),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: http.HandlerFunc(h.PatchTask),
		},
	}
}
