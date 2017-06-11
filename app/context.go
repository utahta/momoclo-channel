package app

import (
	"net/http"

	gcontext "github.com/gorilla/context"
	"github.com/utahta/momoclo-channel/lib/log"
	"golang.org/x/net/context"
)

type Context struct {
	context.Context

	Writer  http.ResponseWriter
	Request *http.Request
}

type contextKey int

const (
	contextVar contextKey = iota
)

func SetContext(req *http.Request, ctx *Context) {
	gcontext.Set(req, contextVar, ctx)
}

func GetContext(req *http.Request) *Context {
	return gcontext.Get(req, contextVar).(*Context)
}

func DeleteContext(req *http.Request) {
	gcontext.Delete(req, contextVar)
}

func (c *Context) Fail(err error) {
	c.Error(err, http.StatusInternalServerError)
}

func (c *Context) Error(err error, code int) {
	var message string
	if err != nil {
		message = err.Error()
	}

	http.Error(c.Writer, message, code)
	log.Errorf(c, "Internal server error! err:%+v", err)
}
