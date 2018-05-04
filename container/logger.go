package container

import (
	"context"

	"github.com/utahta/momoclo-channel/log"
)

// LoggerContainer dependency injection
type LoggerContainer struct {
	ctx context.Context
}

// Logger returns container of logger
func Logger(ctx context.Context) *LoggerContainer {
	return &LoggerContainer{ctx}
}

// AE returns app engine logger
func (c *LoggerContainer) AE() log.Logger {
	return log.NewAELogger(c.ctx)
}
