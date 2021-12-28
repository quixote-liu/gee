package render

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSON struct {
	Data interface{}
}

var (
	jsonContentType      = []string{"application/json; charset=utf-8"}
	jsonpContentType     = []string{"application/javascript; charset=utf-8"}
	jsonASCIIContentType = []string{"application/json"}
)

func (r JSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	bytes, err := json.Marshal(r.Data)
	if err != nil {
		return fmt.Errorf("marshal data failed: %v", err)
	}

	_, err = w.Write(bytes)
	return err
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
