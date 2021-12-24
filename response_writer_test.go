package gee

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	bodyCase = map[string]interface{}{
		"name":    "lcs",
		"age":     "18",
		"address": "hubeidaxue",
	}
)

func TestWrite(t *testing.T) {
	jsonBytes, err := json.Marshal(bodyCase)
	assert.Nil(t, err)

	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var nrw ResponseWriter = &responseWriter{
			ResponseWriter: rw,
		}
		nrw.WriteHeader(http.StatusCreated)
		nrw.Write(jsonBytes)
	}))
	defer s.Close()

	resp, err := http.Get(s.URL)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	payload := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&payload)
	assert.Nil(t, err)

	assert.Equal(t, bodyCase, payload)
}
