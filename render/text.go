package render

import (
	"fmt"
	"gee/internal/bytesconv"
	"net/http"
)

// String contains the given interface object slice and its format.
type String struct {
	Format string
	Data   []interface{}
}

var plainContentType = []string{"text/plain; charset=utf-8"}

// Render (String) writes data with custom ContentType.
func (r String) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	if len(r.Data) > 0 {
		_, err = fmt.Fprintf(w, r.Format, r.Data...)
		return
	}
	_, err = w.Write(bytesconv.StringToBytes(r.Format))
	return
}

// WriteContentType (String) writes Plain ContentType.
func (r String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, plainContentType)
}
