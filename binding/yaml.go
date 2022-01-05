package binding

import (
	"bytes"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (yamlBinding) Bind(w *http.Request, obj interface{}) error {
	return decodeYaml(w.Body, obj)
}

func (yamlBinding) BindBody(body []byte, obj interface{}) error {
	return decodeYaml(bytes.NewReader(body), obj)
}

func decodeYaml(r io.Reader, obj interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	// TODO: Add validator
	return nil
}
