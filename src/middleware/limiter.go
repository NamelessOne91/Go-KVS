package middleware

import "net/http"

const (
	maxSize int64 = 1024
)

// Limiter is an HTTP middleware that wraps a Request's Body
// in a MaxBytesReader allowing to limit the size of incoming clients'
// requests, sent either by accident or with malicious intents
func Limiter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
