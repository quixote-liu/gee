package gee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRouter(http.MethodGet, pattern, handler)
}

func (e *Engine) PUT(pattern string, handler HandlerFunc) {
	e.router.addRouter(http.MethodPut, pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRouter(http.MethodPost, pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
