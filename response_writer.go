package middlewares

import (
	"io"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter

	Status() int
	Size() int
	Written() bool
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	nrw := &responseWriter{
		ResponseWriter: rw,
	}

	return nrw
}

type responseWriter struct {
	http.ResponseWriter

	status int
	size   int
}

func (rw *responseWriter) WriteHeader(s int) {
	if rw.Written() {
		return
	}

	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.size += size

	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) ReadFrom(r io.Reader) (int64, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}

	n, err := io.Copy(rw.ResponseWriter, r)
	rw.size += int(n)

	return n, err
}

// implement http.ResponseController
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
