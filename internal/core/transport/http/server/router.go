package core_http_server

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

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
	methodsByPath := map[string]map[string]struct{}{}

	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.WithMiddleware())

		if methodsByPath[route.Path] == nil {
			methodsByPath[route.Path] = map[string]struct{}{}
		}
		methodsByPath[route.Path][route.Method] = struct{}{}

		if route.Method == http.MethodGet {
			headPattern := fmt.Sprintf("%s %s", http.MethodHead, route.Path)
			r.Handle(headPattern, route.WithMiddleware())
			methodsByPath[route.Path][http.MethodHead] = struct{}{}
		}
	}

	for path, methods := range methodsByPath {
		allowMethods := []string{http.MethodOptions}
		for m := range methods {
			allowMethods = append(allowMethods, m)
		}
		sort.Strings(allowMethods)
		allowedHeader := strings.Join(allowMethods, ", ")

		optionsPattern := fmt.Sprintf("%s %s", http.MethodOptions, path)
		// capture for closure
		allowed := allowedHeader
		r.Handle(optionsPattern, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Allow", allowed)
			w.Header().Set("Access-Control-Allow-Methods", allowed)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusNoContent)
		}))
	}
}

func (r *APIVersionRouters) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
