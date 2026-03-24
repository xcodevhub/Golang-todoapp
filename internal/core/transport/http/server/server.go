package core_http_server

import (
	"context"
	"fmt"
	"net/http"

	core_logger "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/logger"
	core_http_middleware "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func (h *HTTPServer) RegisterAPIRouter(router *APIVersionRouters) {
	h.RegisterHandler(router)
}

// Створює новий сервер
func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

// Реєстрація router-ів з префіксом версії API
func (h *HTTPServer) RegisterHandler(routers ...*APIVersionRouters) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersionRouter)
		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

// Повертає ServeMux для реєстрації handler-ів
func (h *HTTPServer) Mux() *http.ServeMux {
	return h.mux
}

// Run запускає HTTP сервер і чекає на контекст для завершення
func (h *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	// Запускаємо сервер в горутині
	go func() {
		h.log.Warn("starting HTTP server", zap.String("addr", h.config.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- err
		}
		close(ch)
	}()

	// Чекаємо на помилку або сигнал завершення
	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server ...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		h.log.Warn("HTTP server stopped")
	}

	return nil
}
