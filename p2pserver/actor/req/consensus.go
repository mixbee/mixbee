

package req

import (
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
)

var ConsensusPid *actor.PID

func SetConsensusPid(conPid *actor.PID) {
	ConsensusPid = conPid
}
