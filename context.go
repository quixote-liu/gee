package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin object
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		Params: make(map[string]string),
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) setHeader(key, value string) {
	c.Writer.Header().Add(key, value)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.setHeader("Content-Type", "text/plain")
	c.status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.setHeader("Content-Type", "application/json")
	c.status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(&obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) HTML(code int, html string) {
	c.setHeader("Content-Type", "text/html")
	c.status(code)
	c.Writer.Write([]byte(html))
}
