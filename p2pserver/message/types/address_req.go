

package types

import (
	"github.com/mixbee/mixbee/p2pserver/common"
)

type AddrReq struct{}

//Serialize message payload
func (this AddrReq) Serialization() ([]byte, error) {
	return nil, nil
}

func (this *AddrReq) CmdType() string {
	return common.GetADDR_TYPE
}

//Deserialize message payload
func (this *AddrReq) Deserialization(p []byte) error {
	return nil
}
