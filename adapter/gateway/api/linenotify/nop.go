package linenotify

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

type nop struct{}

func (c *nop) Notify(_ string, msg model.LineNotifyMessage) error {
	return nil
}
