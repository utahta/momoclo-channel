package middleware

import (
	"net/http"
	//"time"

	gcontext "github.com/gorilla/context"
	//"golang.org/x/net/context"
	"google.golang.org/appengine"
)

const appengineContextKey string = "appengine-context"

// Appengine middleware
func Appengine(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	//ctx, cancel := context.WithTimeout(appengine.NewContext(req), 60*time.Second)
	//defer cancel()
	ctx := appengine.NewContext(req)

	gcontext.Set(req, appengineContextKey, ctx)

	next(rw, req)

	gcontext.Delete(req, appengineContextKey)
}
