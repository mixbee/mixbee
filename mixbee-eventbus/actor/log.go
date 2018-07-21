package actor

import (
	"github.com/mixbee/mixbee/mixbee-eventbus/log"
)

var (
	plog = log.New(log.DebugLevel, "[ACTOR]")
)

// SetLogLevel sets the log level for the logger.
//
// SetLogLevel is safe to call concurrently
func SetLogLevel(level log.Level) {
	plog.SetLevel(level)
}
