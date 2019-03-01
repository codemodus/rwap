package rwap

import (
	"net/http"
	"strconv"
)

// Rwap ...
type Rwap struct {
	http.ResponseWriter
	status        int
	contentLength int64
}

// New ...
func New(w http.ResponseWriter) *Rwap {
	return &Rwap{ResponseWriter: w}
}

func (w *Rwap) Write(p []byte) (int, error) {
	cl, err := w.ResponseWriter.Write(p)
	w.contentLength += int64(cl)

	return cl, err
}

// WriteHeader ...
func (w *Rwap) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Status ...
func (w *Rwap) Status() int {
	return w.status
}

// ContentLength ...
func (w *Rwap) ContentLength() int64 {
	if w.contentLength > 0 {
		return w.contentLength
	}

	cl, err := strconv.ParseInt(w.Header().Get("Content-Length"), 10, 64)
	if err != nil {
		return 0
	}

	return cl
}

// Wrap ...
func Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(New(w), r)
	})
}
