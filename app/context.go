package app

import (
	"net/http"

	gcontext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

type Context struct {
	context.Context
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
