package model

import (
	"github.com/mjibson/goon"
	"golang.org/x/net/context"
)

type UstreamStatus struct {
	Id   string `datastore:"-" goon:"id"`
	IsLive bool
}

func NewUstreamStatus() *UstreamStatus {
	return &UstreamStatus{
		Id:   "ustream_status",
		IsLive: false,
	}
}

func (u *UstreamStatus) Get(ctx context.Context) {
	g := goon.FromContext(ctx)
	g.Get(u)
}

func (u *UstreamStatus) Put(ctx context.Context) {
	g := goon.FromContext(ctx)
	g.Put(u)
}
