

package types

import (
	"bytes"
	"fmt"

	ct "github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
)

type Block struct {
	Blk ct.Block
}

//Serialize message payload
func (this Block) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := this.Blk.Serialize(p)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. Blk:%v", this.Blk))
	}

	return p.Bytes(), nil
}

func (this *Block) CmdType() string {
	return common.BLOCK_TYPE
}

//Deserialize message payload
func (this *Block) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := this.Blk.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read Blk error. buf:%v", buf))
	}

	return nil
}
