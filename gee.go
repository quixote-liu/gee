package gee

import (
	"html/template"
	"net/http"
	"strings"
)

const defaultMultipartMemory = 32 << 20 // 32MB

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups

	// for html render
	htmlTemplates    *template.Template
	funcMap          template.FuncMap
	secureJSONPrefix string

	// Value of "maxMemory" param that is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64
}

func New() *Engine {
	engine := &Engine{
		router:             newRouter(),
		secureJSONPrefix:   "while(1);",
		MaxMultipartMemory: defaultMultipartMemory,
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) SecureJsonPrefix(prefix string) *Engine {
	e.secureJSONPrefix = prefix
	return e
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
	// combine handlers by group prefix
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.handlers...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = e
	e.router.handle(c)
}

// start server
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(
		template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}
