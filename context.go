package gee

import (
	"gee/render"
	"net/http"
	"net/url"
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

	queryCache url.Values
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

/**************************************************/
/****************** INPUT DATA ********************/
/**************************************************/

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) setHeader(key, value string) {
	c.Writer.Header().Add(key, value)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.Req != nil {
			c.queryCache = c.Req.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

// Query returns the keyed url query value if is exist,
// Otherwise it returns an empty string `("")`.
func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

// GetQuery is like Query, it returns the keyed url query value
// if it exist `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`
func (c *Context) GetQuery(key string) (value string, ok bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

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

func (c *Context) JSON(code int, obj interface{}) {
	// do not need the pointer of render.JSON,
	// the render.JSON structure is main to compose the render work.
	c.Render(code, render.JSON{Data: obj})
}

func (c *Context) IndentedJSON(code int, obj interface{}) {
	c.Render(code, render.IndentedJSON{Data: obj})
}

func (c *Context) SecureJSON(code int, obj interface{}) {
	c.Render(code, render.SecureJSON{Data: obj, Prefix: "prefix"})
}

func (c *Context) JSONP(code int, obj interface{}) {
	callback := c.DefaultQuery("callback", "")
	if callback == "" {
		c.Render(code, render.JSON{Data: obj})
		return
	}
	c.Render(code, render.JsonpJSON{Callback: callback, Data: obj})
}

func (c *Context) AsciiJSON(code int, obj interface{}) {
	c.Render(code, render.AsciiJSON{Data: obj})
}

func (c *Context) XML(code int, obj interface{}) {
	c.Render(code, render.XML{Data: obj})
}

func (c *Context) YAML(code int, obj interface{}) {
	c.Render(code, render.YAML{Data: obj})
}

func (c *Context) String(code int, format string, obj ...interface{}) {
	c.Render(code, render.String{Format: format, Data: obj})
}

func (c *Context) Fail(code int, err error) {
	c.Status(code)
	c.Writer.Write([]byte(err.Error()))
}

func (c *Context) HTML(code int, name string, data interface{}) {
	// c.setHeader("Content-Type", "text/html")
	// c.Status(code)
	// if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
	// 	c.Fail(code, err)
	// }
	c.Render(code, render.HTML{Template: c.engine.htmlTemplates, Name: name, Data: data})
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}
