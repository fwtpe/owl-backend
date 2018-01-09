package client

import (
	log "github.com/fwtpe/owl/common/logruslog"
)

var logger = log.NewDefaultLogger("INFO")

func SetLoggerLevel(level string) {
	logger = log.NewDefaultLogger(level)
}
