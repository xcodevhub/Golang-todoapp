package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/logger"
	core_pgx_pool "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/transport/http/middleware"
	core_http_server "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/features/tasks/repository/postgres"
	tasks_service "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/features/tasks/service"
	tasks_transport_http "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/features/tasks/transport/http"
	users_repository_postgres "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/features/users/repository/postgres"
	users_service "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/features/users/service"
	users_transport_http "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Printf("failed to init application logger: %v", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Debug("application time zone ", zap.Any("zone", timeZone))

	logger.Debug("Initializing postgres connection pool!")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres conection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("Initializing features", zap.String("feature", "users"))
	usersRepository := users_repository_postgres.NewUsersPostgresRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("Initializing features", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasks, err := tasksService.GetTasks(ctx, nil, nil, nil)
	if err != nil {
		logger.Error("failed to get tasks", zap.Error(err))
	} else {
		logger.Debug("successfully fetched tasks during startup", zap.Int("count", len(tasks)))
	}
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("Initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.LoggerMiddleware(logger),
		core_http_middleware.Trace(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRouters(core_http_server.ApiVersionRouter1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouter(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
