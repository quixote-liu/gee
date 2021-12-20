package gee

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePattern(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(parsePattern("/name/:id"), []string{"name", ":id"})
	assert.Equal(parsePattern("/name/*"), []string{"name", "*"})
	assert.Equal(parsePattern("/name/*id/*"), []string{"name", "*id"})
}

func TestGetRouter(t *testing.T) {
	r := newRouter()
	assert := assert.New(t)
	getMethod := http.MethodGet

	pathCaseA := "/identity/v3/user"
	pathCaseB := "/identity/v3/user/:id"
	pathCaseC := "/identity/v3/user/:id/group"
	pathCaseD := "/identity/v3/*path"

	// add router
	r.addRouter(getMethod, pathCaseA, nil)
	r.addRouter(getMethod, pathCaseB, nil)
	r.addRouter(getMethod, pathCaseC, nil)
	r.addRouter(getMethod, pathCaseD, nil)

	// get router
	n, params := r.getRouter(getMethod, "/identity/v3/user")
	assert.NotNil(n)
	assert.Equal(pathCaseA, n.pattern)
	assert.Equal(0, len(params))

	n, params = r.getRouter(getMethod, "/identity/v3/user/10086")
	assert.NotNil(n)
	assert.Equal(pathCaseB, n.pattern)
	assert.Equal("10086", params["id"])

	n, params = r.getRouter(getMethod, "/identity/v3/user/10086/group")
	assert.NotNil(n)
	assert.Equal(pathCaseC, n.pattern)
	assert.Equal("10086", params["id"])

	n, params = r.getRouter(getMethod, "/identity/v3/helloAllPath/paths")
	assert.NotNil(n)
	assert.Equal(pathCaseD, n.pattern)
	assert.Equal("helloAllPath/paths", params["path"])
}
