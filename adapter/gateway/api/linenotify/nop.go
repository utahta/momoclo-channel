package linenotify

import (
	"github.com/utahta/momoclo-channel/types"
)

type nop struct{}

// NewNop returns no operation LineNotify
func NewNop() types.LineNotify {
	return &nop{}
}

func (c *nop) Notify(_ string, msg types.LineNotifyMessage) error {
	return nil
}
