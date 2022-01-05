package binding

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var JSONdelivery = map[string]interface{}{
	"name":    "lcs",
	"id":      "20202111",
	"email":   "123@foxmail.com",
	"address": "hubeiwuhan",
}

func bodyCase() ([]byte, error) {
	bytes, err := json.Marshal(JSONdelivery)
	return bytes, err
}

func TestBind(t *testing.T) {
	body, err := bodyCase()
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodGet, "/jsonbinding/test", bytes.NewReader(body))
	assert.Nil(t, err)

	jsonBinding := &jsonBinding{}
	payload := map[string]interface{}{}
	err = jsonBinding.Bind(req, &payload)
	assert.Nil(t, err)

	assert.Equal(t, JSONdelivery, payload)
}

func TestBindBody(t *testing.T) {
	body, err := bodyCase()
	assert.Nil(t, err)

	jsonBinding := &jsonBinding{}
	payload := map[string]interface{}{}
	err = jsonBinding.BindBody(body, &payload)
	assert.Nil(t, err)

	assert.Equal(t, JSONdelivery, payload)
}