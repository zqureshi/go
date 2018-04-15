package log

import (
	log "github.com/sirupsen/logrus"
)

var (
	// Logger should be used by every package in repo.
	Logger = log.New()
)

// SetLogger updates the shared root logger.
func SetLogger(l *log.Logger) {
	Logger = l
}
