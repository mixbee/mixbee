

package types

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mixbee/mixbee/common/serialization"
	ct "github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
)

type BlkHeader struct {
	BlkHdr []*ct.Header
}

//Serialize message payload
func (this BlkHeader) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	serialization.WriteUint32(p, uint32(len(this.BlkHdr)))
	for _, header := range this.BlkHdr {
		err := header.Serialize(p)
		if err != nil {
			return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. header:%v", header))
		}
	}

	return p.Bytes(), nil
}

func (this *BlkHeader) CmdType() string {
	return common.HEADERS_TYPE
}

//Deserialize message payload
func (this *BlkHeader) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	var count uint32

	err := binary.Read(buf, binary.LittleEndian, &count)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read Cnt error. buf:%v", buf))
	}

	for i := 0; i < int(count); i++ {
		var headers ct.Header
		err := headers.Deserialize(buf)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("deserialize headers error. buf:%v", buf))
		}
		this.BlkHdr = append(this.BlkHdr, &headers)
	}
	return nil
}
