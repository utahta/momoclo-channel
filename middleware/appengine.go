package middleware

import (
	"net/http"

	"github.com/utahta/momoclo-channel/app"
	"google.golang.org/appengine"
)

// Appengine middleware
func Appengine(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ctx := &app.Context{
		Context: appengine.NewContext(req),
		Writer:  rw,
		Request: req,
	}
	app.SetContext(req, ctx)
	defer app.DeleteContext(req)

	next(rw, req)
}
