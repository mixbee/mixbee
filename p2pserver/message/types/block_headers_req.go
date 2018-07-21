

package types

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
)

type HeadersReq struct {
	Len       uint8
	HashStart [common.HASH_LEN]byte
	HashEnd   [common.HASH_LEN]byte
}

//Serialize message payload
func (this *HeadersReq) Serialization() ([]byte, error) {
	p := new(bytes.Buffer)
	err := binary.Write(p, binary.LittleEndian, this)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("binary.Write payload error. payload:%v", this))
	}

	return p.Bytes(), nil
}

func (this *HeadersReq) CmdType() string {
	return common.GET_HEADERS_TYPE
}

//Deserialize message payload
func (this *HeadersReq) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, this)

	return err
}
