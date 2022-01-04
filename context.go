package gee

import (
	"errors"
	"gee/binding"
	"gee/render"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
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

	// This mutex protect Keys map.
	mu sync.RWMutex

	// Keys is key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	// middleware
	handlers []HandlerFunc
	index    int

	// engine pointer
	engine *Engine

	// queryCache store the requested all query string.
	queryCache url.Values

	// formCache use url.ParseQuery cached postForm contains the parsed form data from POST, PATCH,
	// or PUT body parameters.
	formCache url.Values
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

/**********************************************************/
/******************* METADATA MANAGEMENT ******************/
/**********************************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes c.Keys if it was not used proviously.
func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true)
// If the value does not exist it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// MustGet returns the value for the given key if it exist, otherwise it panics.
func (c *Context) MustGet(key string) (value interface{}) {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

/*****************************************************/
/********************* INPUT DATA ********************/
/*****************************************************/

// Param returns the value of the URL param.
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// AddParam adds param to context and
// replaces path param key with given value for e2e testing purpose.
func (c *Context) AddParam(key, value string) {
	c.Params[key] = value
}

// Query returns the keyed url query value if it exists,
// Otherwise it returns an empty string `("")`.
func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
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

// QueryArray returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) QueryArray(key string) (values []string) {
	values, _ = c.GetQueryArray(key)
	return
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

// GetQueryArray returns the value of associated key,
func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

// PostForm returns the specified key from a POST urlencoded form or multipart form
// when is exists, otherwise it returns an empty value.
func (c *Context) PostForm(key string) string {
	val, _ := c.GetPostForm(key)
	return val
}

// GetPostForm returns a string of a given form key
func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// PostFormArray returns the a slice of strings for a given form key.
func (c *Context) PostFormArray(key string) (values []string) {
	values, _ = c.GetPostFormArray(key)
	return values
}

func (c *Context) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		req := c.Req
		if err := req.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
			if errors.Is(err, http.ErrNotMultipart) {
				debugPrint("error on parse multipart form array: %v", err)
			}
		}
		c.formCache = req.PostForm
	}
}

// GetPostFormArray returns a slice of strings of a given form key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetPostFormArray(key string) (values []string, ok bool) {
	c.initFormCache()
	values, ok = c.formCache[key]
	return
}

// ShouldBindWith binds the http passed struct pointer using the specified binding engine.
// See the binding package.
func (c *Context) shouldBindWith(obj interface{}, b binding.Binding) error {
	return b.Bind(c.Req, obj)
}

// ShouldBindJSON is a shortcut for c.shouldBindWith(obj, binding.JSON).
func (c *Context) ShouldBindJSON(obj interface{}) error {
	return c.shouldBindWith(obj, binding.JSON)
}

// Bind checks the Content-Type to select a binding engine automatically.
func (c *Context) Bind(obj interface{}) error {
	return nil
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Req.ParseMultipartForm(c.engine.MaxMultipartMemory)
	return c.Req.MultipartForm, err
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
	c.Render(code, render.HTML{Template: c.engine.htmlTemplates, Name: name, Data: data})
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}
