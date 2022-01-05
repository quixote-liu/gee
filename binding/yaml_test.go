package binding

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestYamlBind(t *testing.T) {
	yamlBinding := &yamlBinding{}
	delivery := map[string]interface{}{
		"name":    "lcs",
		"id":      "88888888",
		"address": "hubeiwuhan",
	}
	body, err := yaml.Marshal(delivery)
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "/target", bytes.NewReader(body))

	payload := map[string]interface{}{}
	err = yamlBinding.Bind(req, payload)
	assert.Nil(t, err)
}

func TestYamlBindBody(t *testing.T) {
	yamlBinding := &yamlBinding{}
	delivery := map[string]interface{}{
		"name":    "lcs",
		"id":      "88888888",
		"address": "hubeiwuhan",
	}
	body, err := yaml.Marshal(delivery)
	assert.Nil(t, err)

	payload := map[string]interface{}{}
	err = yamlBinding.BindBody(body, payload)
	assert.Nil(t, err)
}
