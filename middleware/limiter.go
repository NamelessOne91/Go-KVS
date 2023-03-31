package middleware

import "net/http"

const (
	maxSize int64 = 1024
)

func Limiter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
