

package types

import (
	"bytes"
	"fmt"

	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
)

type VerACK struct {
	IsConsensus bool
}

//Serialize message payload
func (this VerACK) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := serialization.WriteBool(p, this.IsConsensus)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. IsConsensus:%v", this.IsConsensus))
	}

	return p.Bytes(), nil
}

func (this VerACK) CmdType() string {
	return common.VERACK_TYPE
}

//Deserialize message payload
func (this *VerACK) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)

	isConsensus, err := serialization.ReadBool(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read IsConsensus error. buf:%v", buf))
	}

	this.IsConsensus = isConsensus
	return nil
}
