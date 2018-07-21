

package states

import (
	"io"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
)

// Invoke smart contract struct
// Param Version: invoke smart contract version, default 0
// Param Address: invoke on blockchain smart contract by address
// Param Method: invoke smart contract method, default ""
// Param Args: invoke smart contract arguments
type Contract struct {
	Version byte
	Address common.Address
	Method  string
	Args    []byte
}

// Serialize contract
func (this *Contract) Serialize(w io.Writer) error {
	if err := serialization.WriteByte(w, this.Version); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Version serialize error!")
	}
	if err := this.Address.Serialize(w); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Address serialize error!")
	}
	if err := serialization.WriteVarBytes(w, []byte(this.Method)); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Method serialize error!")
	}
	if err := serialization.WriteVarBytes(w, this.Args); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Args serialize error!")
	}
	return nil
}

// Deserialize contract
func (this *Contract) Deserialize(r io.Reader) error {
	var err error
	this.Version, err = serialization.ReadByte(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Version deserialize error!")
	}

	if err := this.Address.Deserialize(r); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Address deserialize error!")
	}

	method, err := serialization.ReadVarBytes(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Method deserialize error!")
	}
	this.Method = string(method)

	this.Args, err = serialization.ReadVarBytes(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Args deserialize error!")
	}
	return nil
}

type PreExecResult struct {
	State  byte
	Gas    uint64
	Result interface{}
}
