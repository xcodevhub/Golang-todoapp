package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/logger"
	core_http_request "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/transport/http/request"
	core_http_response "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'user_id'/'limit'/'offset' query params",
		)
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(tasksDTOsFromDomains(tasksDomains))

	responseHandler.JsonResponse(response, http.StatusOK)

}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		UserIDQueryParamKey = "user_id"
		LimitQueryParamKey  = "limit"
		OffsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, UserIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, LimitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get limit query param: %w", err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, OffsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	return userID, limit, offset, nil

}
