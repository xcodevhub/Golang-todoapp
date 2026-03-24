package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/transport/http/middleware"
)

type Route struct {
	Method      string                            // POST, GET, etc.
	Path        string                            // /users, /users/{id}, etc.
	Handler     http.Handler                      // http.HandlerFunc(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { ... }))
	middlewares []core_http_middleware.Middleware // middlewares to be applied to this route
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.middlewares...,
	)

}
