

package types

import (
	"bytes"
	"fmt"

	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
)

type Consensus struct {
	Cons ConsensusPayload
}

//Serialize message payload
func (this *Consensus) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := this.Cons.Serialize(p)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. consensus:%v", this.Cons))
	}
	return p.Bytes(), nil
}

func (this *Consensus) CmdType() string {
	return common.CONSENSUS_TYPE
}

//Deserialize message payload
func (this *Consensus) Deserialization(p []byte) error {
	log.Debug()
	buf := bytes.NewBuffer(p)
	err := this.Cons.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("deserialize Node error. buf:%v", buf))
	}
	return nil
}
