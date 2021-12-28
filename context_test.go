package gee

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test/contextJSON", nil)

	context := newContext(w, r)
	data := H{
		"name":    "lcs",
		"id":      "20202111",
		"address": "hubeidaxue",
	}
	context.JSON(http.StatusOK, data)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	actual := H{}
	err := json.NewDecoder(w.Body).Decode(&actual)
	assert.Nil(t, err)
	assert.Equal(t, data, actual)
}
