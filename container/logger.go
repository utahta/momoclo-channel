package container

import (
	"github.com/utahta/momoclo-channel/log"
)

// LoggerContainer dependency injection
type LoggerContainer struct {
}

// Logger returns container of logger
func Logger() *LoggerContainer {
	return &LoggerContainer{}
}

// AE returns app engine logger
func (c *LoggerContainer) AE() log.Logger {
	return log.NewAELogger()
}
