package core_http_server

import "net/http"

type Route struct {
	Method  string       // POST, GET, etc.
	Path    string       // /users, /users/{id}, etc.
	Handler http.Handler // http.HandlerFunc(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { ... }))
}

func NewRoute(
	method string,
	path string,
	handler http.Handler,
) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
