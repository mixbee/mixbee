

package req

import (
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
)

var CrossChainPid *actor.PID

func SetCrossChainPid(Pid *actor.PID) {
	CrossChainPid = Pid
}
