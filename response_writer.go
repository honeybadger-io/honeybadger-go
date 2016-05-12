package honeybadger

import "net/http"

type responseWriter struct {
	status int
	http.ResponseWriter
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{http.StatusOK, w}
}
