

package types

import (
	"bytes"
	"fmt"

	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
)

type Ping struct {
	Height uint64
}

//Serialize message payload
func (this Ping) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := serialization.WriteUint64(p, this.Height)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. Height:%v", this.Height))
	}

	return p.Bytes(), nil
}

func (this *Ping) CmdType() string {
	return common.PING_TYPE
}

//Deserialize message payload
func (this *Ping) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)

	height, err := serialization.ReadUint64(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read Height error. buf:%v", buf))
	}
	this.Height = height
	return nil
}
