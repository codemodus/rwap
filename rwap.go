package rwap

import (
	"net/http"
	"strconv"
)

type ResponseWriter struct {
	http.ResponseWriter
	status        int
	contentLength int64
}

func (w *ResponseWriter) Write(p []byte) (int, error) {
	cl, err := w.ResponseWriter.Write(p)
	w.contentLength += int64(cl)

	return cl, err
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseWriter) Status() int {
	return w.status
}

func (w *ResponseWriter) ContentLength() int64 {
	if w.contentLength > 0 {
		return w.contentLength
	}

	cl, err := strconv.ParseInt(w.Header().Get("Content-Length"), 10, 64)
	if err != nil {
		return 0
	}

	return cl
}

func Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := &ResponseWriter{ResponseWriter: w}
		next.ServeHTTP(nw, r)
	})
}
