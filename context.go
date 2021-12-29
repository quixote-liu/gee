package gee

import (
	"fmt"
	"gee/render"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin object
	Writer ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int

	// middleware
	handlers []HandlerFunc
	index    int

	// engine pointer
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	nw := &responseWriter{}
	nw.reset(w)
	return &Context{
		Writer: nw,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		Params: make(map[string]string),
		index:  -1,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) setHeader(key, value string) {
	c.Writer.Header().Add(key, value)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

// func (c *Context) GetQuery(key string) (value string, ok bool) {
// 	if c.Req.URL.Query().Get(key)
// }

// func (c *Context) DefaultQuery(key, defaultValue string) string {

// }

/**************************************************/
/************ RESPONSE RENDERING ******************/
/**************************************************/

// Status set the HTTP response status
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func (c *Context) bodyAllowedForStatus(status int) bool {
	switch {
	case status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !c.bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	}

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.setHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	// do not need the pointer of render.JSON,
	// the render.JSON structure is main to compose the render work.
	c.Render(code, render.JSON{Data: obj})
}

// func (c *Context) JSONP(code int, obj interface{}) {
// 	c.Render(code, render.JsonpJSON{Data: obj})
// }

func (c *Context) Fail(code int, err error) {
	c.Status(code)
	c.Writer.Write([]byte(err.Error()))
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.setHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(code, err)
	}
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}
