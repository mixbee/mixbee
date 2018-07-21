

package types

import (
	"bytes"
	"fmt"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/errors"
	common2 "github.com/mixbee/mixbee/p2pserver/common"
)

type NotFound struct {
	Hash common.Uint256
}

//Serialize message payload
func (this NotFound) Serialization() ([]byte, error) {
	return this.Hash[:], nil
}

func (this NotFound) CmdType() string {
	return common2.NOT_FOUND_TYPE
}

//Deserialize message payload
func (this *NotFound) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)

	err := this.Hash.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("deserialize Hash error. buf:%v", buf))
	}
	return nil
}
