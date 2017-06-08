package middleware

import (
	"net/http"

	gcontext "github.com/gorilla/context"
	"google.golang.org/appengine"
)

const appengineContextKey string = "appengine-context"

// Appengine middleware
func Appengine(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ctx := appengine.NewContext(req)
	gcontext.Set(req, appengineContextKey, ctx)
	defer gcontext.Delete(req, appengineContextKey)

	next(rw, req)
}
