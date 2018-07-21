

package types

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
)

type BlocksReq struct {
	HeaderHashCount uint8
	HashStart       [common.HASH_LEN]byte
	HashStop        [common.HASH_LEN]byte
}

//Serialize message payload
func (this *BlocksReq) Serialization() ([]byte, error) {
	p := new(bytes.Buffer)
	err := binary.Write(p, binary.LittleEndian, this)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write  error. payload:%v", this))
	}

	return p.Bytes(), nil
}

func (this *BlocksReq) CmdType() string {
	return common.GET_BLOCKS_TYPE
}

//Deserialize message payload
func (this *BlocksReq) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, this)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read BlocksReq error. buf:%v", buf))
	}
	return nil
}
