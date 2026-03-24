package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/transport/http/middleware"
)

type ApiVersionRouter string

var (
	ApiVersionRouter1 ApiVersionRouter = ("v1")
	ApiVersionRouter2 ApiVersionRouter = ("v2")
	ApiVersionRouter3 ApiVersionRouter = ("v3")
)

type APIVersionRouters struct {
	*http.ServeMux
	apiVersionRouter ApiVersionRouter
	middleware       []core_http_middleware.Middleware
}

func NewAPIVersionRouters(
	apiVersionRouter ApiVersionRouter,
	middleware ...core_http_middleware.Middleware,
) *APIVersionRouters {
	return &APIVersionRouters{
		ServeMux:         http.NewServeMux(),
		apiVersionRouter: apiVersionRouter,
		middleware:       middleware,
	}
}

func (r *APIVersionRouters) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		// "GET /users" etc.
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouters) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
