package render

import (
	"net/http"
)

// Render interface is to be implemented by JSON, XML, HTML, and so on.
type Render interface {
	// Render writes data with custom ContentType
	Render(w http.ResponseWriter) error
	// WriteContentType writes custom ContentType
	WriteContentType(w http.ResponseWriter)
}

// writeContentType write Content-Type value to header
func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
