

package types

import (
	"github.com/mixbee/mixbee/p2pserver/common"
)

type Disconnected struct{}

//Serialize message payload
func (this Disconnected) Serialization() ([]byte, error) {
	return nil, nil
}

func (this Disconnected) CmdType() string {
	return common.DISCONNECT_TYPE
}

//Deserialize message payload
func (this *Disconnected) Deserialization(p []byte) error {
	return nil
}
