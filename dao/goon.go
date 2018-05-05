package dao

import (
	"context"

	"github.com/mjibson/goon"
)

var ctxKey struct{}

// WithGoon returns context with goon
func WithGoon(ctx context.Context, g *goon.Goon) context.Context {
	return context.WithValue(ctx, ctxKey, g)
}

// FromContext returns goon with this context for key
func FromContext(ctx context.Context) *goon.Goon {
	v := ctx.Value(ctxKey)
	g, ok := v.(*goon.Goon)
	if !ok {
		return goon.FromContext(ctx)
	}
	return g
}
