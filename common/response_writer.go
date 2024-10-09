package common

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// A ResponseWriter wraps an [http.ResponseWriter] to add additional functionality.
type ResponseWriter interface {
	http.ResponseWriter

	Status() int
	Size() int
	Written() bool
}

// NewResponseWriter returns a ResponseWriter wrapping the provided [http.ResponseWriter].
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	nrw := &responseWriter{
		ResponseWriter: rw,
	}

	return nrw
}

// A responseWriter wraps an [http.ResponseWriter] to record the number of bytes written in a response.
type responseWriter struct {
	http.ResponseWriter

	status int
	size   int
}

// Implements [http.ResponseWriter.WriteHeader]
func (rw *responseWriter) WriteHeader(s int) {
	if rw.Written() {
		return
	}

	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

// Implements [http.ResponseWriter.Write]
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.size += size

	return size, err
}

// Status retrieves the recorded response status code.
func (rw *responseWriter) Status() int {
	return rw.status
}

// Size retrieves the recorded response bytes written.
func (rw *responseWriter) Size() int {
	return rw.size
}

// Written retrieves whether or not the response headers have already been written.
func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

// Hijack implements the [http.Hijacker] interface.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if orw, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return orw.Hijack()
	}

	return nil, nil, fmt.Errorf("Hijacker interface not implemented")
}

// Flush implements the [http.Flusher] interface.
func (rw *responseWriter) Flush() {
	if orw, ok := rw.ResponseWriter.(http.Flusher); ok {
		orw.Flush()
	}
}
