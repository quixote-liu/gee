package gee

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher

	// Status returns the http response status code of the current request
	Status() int

	// Returns the number of bytes already written into the response body
	// See written()
	Size() int

	// Writes the string into the response body
	WriteString(string) (int, error)

	// Returns true if the response body was already written
	Written() bool

	// Fouces to write the http header (status code + headers)
	WriteHeaderNow()

	// Get the http.Pusher for server push
	Pusher() http.Pusher
}

type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

var _ ResponseWriter = &responseWriter{}

// reset set responseWriter to default value
func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.size != noWritten
}

func (w *responseWriter) WriterHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			log.Printf("[WARING] Headers were already written, Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

// writer status code
func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

// Write write data to response body.
func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	return
}

// Write string value to responseWriter.
// use io.WriteString is a goog way.
func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size = n
	return
}

// Hijack implements the http.Hijack interface.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// Flush implements the http.Flush interface.
func (w *responseWriter) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

// Flush implements the http.Flush interface.
func (w *responseWriter) Pusher() http.Pusher {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
