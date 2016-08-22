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

func (u *UstreamStatus) Load(ctx context.Context) {
	g := goon.FromContext(ctx)
	g.Get(u)
}

func (u *UstreamStatus) Update(ctx context.Context) {
	g := goon.FromContext(ctx)
	g.Put(u)
}
