package server

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       bytes.Buffer
	written    bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
}

func (w *ResponseWriter) WriteHeader(code int) {
	if !w.written {
		w.StatusCode = code
		w.written = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.written = true
	}
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ResponseWriter) Written() bool {
	return w.written
}

func EnsureResponseWriter(w http.ResponseWriter) *ResponseWriter {
	if rw, ok := w.(*ResponseWriter); ok {
		return rw
	}
	return NewResponseWriter(w)
}
