package middlewares

import "net/http"

// WithContentJSON middleware adds Context-Type JSON to response
func WithContentJSON(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	}
}
