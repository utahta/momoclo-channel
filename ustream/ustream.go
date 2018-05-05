package ustream

import (
	"context"

	"github.com/utahta/ustream-live-status"
	"google.golang.org/appengine/urlfetch"
)

type (
	// StatusChecker interface
	StatusChecker interface {
		IsLive(ctx context.Context) (bool, error)
	}

	statusChecker struct {
	}
)

// NewStatusChecker returns UstreamStatusChecker wraps ustream live status client
func NewStatusChecker() StatusChecker {
	return &statusChecker{}
}

// IsLive returns true if status is live
func (s *statusChecker) IsLive(ctx context.Context) (bool, error) {
	c, err := uststat.New(uststat.WithHTTPTransport(&urlfetch.Transport{Context: ctx}))
	if err != nil {
		return false, err
	}
	return c.IsLiveByChannelID("4979543") // momoclotv
}
