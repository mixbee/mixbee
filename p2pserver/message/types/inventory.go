

package types

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
	p2pCommon "github.com/mixbee/mixbee/p2pserver/common"
)

var LastInvHash common.Uint256

type InvPayload struct {
	InvType common.InventoryType
	Blk     []common.Uint256
}

type Inv struct {
	P InvPayload
}

func (this Inv) invType() common.InventoryType {
	return this.P.InvType
}

func (this *Inv) CmdType() string {
	return p2pCommon.INV_TYPE
}

//Serialize message payload
func (this Inv) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := serialization.WriteUint8(p, uint8(this.P.InvType))
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. InvType:%v", this.P.InvType))
	}
	blkCnt := uint32(len(this.P.Blk))
	err = serialization.WriteUint32(p, blkCnt)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. Cnt:%v", blkCnt))
	}
	err = binary.Write(p, binary.LittleEndian, this.P.Blk)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. Blk:%v", this.P.Blk))
	}

	return p.Bytes(), nil
}

//Deserialize message payload
func (this *Inv) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	invType, err := serialization.ReadUint8(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read invType error. buf:%v", buf))
	}
	this.P.InvType = common.InventoryType(invType)
	blkCnt, err := serialization.ReadUint32(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read Cnt error. buf:%v", buf))
	}
	if blkCnt > p2pCommon.MAX_INV_BLK_CNT {
		blkCnt = p2pCommon.MAX_INV_BLK_CNT
	}
	for i := 0; i < int(blkCnt); i++ {
		var blk common.Uint256
		err := binary.Read(buf, binary.LittleEndian, &blk)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read inv blk error. buf:%v", buf))
		}
		this.P.Blk = append(this.P.Blk, blk)
	}
	return nil
}
