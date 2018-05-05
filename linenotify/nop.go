package linenotify

import "context"

type nop struct{}

// NewNop return no operation client
func NewNop() Client {
	return &nop{}
}

func (c *nop) Notify(_ context.Context, _ string, msg Message) error {
	return nil
}
