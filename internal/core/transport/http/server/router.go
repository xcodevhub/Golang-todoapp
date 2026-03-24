package core_http_server

import (
	"fmt"
	"net/http"
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
}

func NewAPIVersionRouters(
	apiVersionRouter ApiVersionRouter,
) *APIVersionRouters {
	return &APIVersionRouters{
		ServeMux:         http.NewServeMux(),
		apiVersionRouter: apiVersionRouter,
	}
}

func (r *APIVersionRouters) RegisterRoutes(route ...Route) {
	for _, route := range route {
		//"GET/tasks"
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.Handler)

	}
}
