package middleware

import (
	"net/http"

	"google.golang.org/appengine"
)

// Appengine wraps appengine context.
func Appengine(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.WithContext(r.Context(), r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
