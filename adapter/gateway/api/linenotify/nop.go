package linenotify

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

type nop struct{}

// NewNop returns no operation LineNotify
func NewNop() model.LineNotify {
	return &nop{}
}

func (c *nop) Notify(_ string, msg model.LineNotifyMessage) error {
	return nil
}
