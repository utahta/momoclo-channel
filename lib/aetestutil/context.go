package aetestutil

import (
	"context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

// NewContex extends aetest.NewContext
func NewContex(opts *aetest.Options) (context.Context, func(), error) {
	inst, err := aetest.NewInstance(opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		inst.Close()
		return nil, nil, err
	}
	ctx := appengine.NewContext(req)
	return ctx, func() {
		inst.Close()
	}, nil
}
