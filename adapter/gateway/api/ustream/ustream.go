package ustream

import (
	"context"

	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/ustream-live-status"
	"google.golang.org/appengine/urlfetch"
)

type (
	statusChecker struct {
		*uststat.Client
	}
)

// NewStatusChecker returns UstreamStatusChecker wraps ustream live status client
func NewStatusChecker(ctx context.Context) types.UstreamStatusChecker {
	c, _ := uststat.New(uststat.WithHTTPTransport(&urlfetch.Transport{Context: ctx}))
	return &statusChecker{Client: c}
}

// IsLive returns true if status is live
func (s *statusChecker) IsLive() (bool, error) {
	return s.IsLiveByChannelID("4979543") // momoclotv
}
