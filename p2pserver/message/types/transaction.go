

package types

import (
	"bytes"
	"fmt"

	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
)

// Transaction message
type Trn struct {
	Txn *types.Transaction
}

//Serialize message payload
func (this Trn) Serialization() ([]byte, error) {
	p := bytes.NewBuffer(nil)
	err := this.Txn.Serialize(p)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. Txn:%v", this.Txn))
	}

	return p.Bytes(), nil
}

func (this *Trn) CmdType() string {
	return common.TX_TYPE
}

//Deserialize message payload
func (this *Trn) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)

	tx := types.Transaction{}
	err := tx.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read txn error. buf:%v", buf))
	}

	this.Txn = &tx
	return nil
}
