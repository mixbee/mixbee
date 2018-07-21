

package ont

import (
	"fmt"
	"io"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
)

// Transfers
type Transfers struct {
	States []*State
}

func (this *Transfers) Serialize(w io.Writer) error {
	if err := utils.WriteVarUint(w, uint64(len(this.States))); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TokenTransfer] Serialize States length error!")
	}
	for _, v := range this.States {
		if err := v.Serialize(w); err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[TokenTransfer] Serialize States error!")
		}
	}
	return nil
}

func (this *Transfers) Deserialize(r io.Reader) error {
	n, err := utils.ReadVarUint(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TokenTransfer] Deserialize states length error!")
	}
	for i := 0; uint64(i) < n; i++ {
		state := new(State)
		if err := state.Deserialize(r); err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[TokenTransfer] Deserialize states error!")
		}
		this.States = append(this.States, state)
	}
	return nil
}

type State struct {
	From  common.Address
	To    common.Address
	Value uint64
}

func (this *State) Serialize(w io.Writer) error {
	if err := utils.WriteAddress(w, this.From); err != nil {
		return fmt.Errorf("[State] serialize from error:%v", err)
	}
	if err := utils.WriteAddress(w, this.To); err != nil {
		return fmt.Errorf("[State] serialize to error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.Value); err != nil {
		return fmt.Errorf("[State] serialize value error:%v", err)
	}
	return nil
}

func (this *State) Deserialize(r io.Reader) error {
	var err error
	this.From, err = utils.ReadAddress(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize from error:%v", err)
	}
	this.To, err = utils.ReadAddress(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize to error:%v", err)
	}

	this.Value, err = utils.ReadVarUint(r)
	if err != nil {
		return err
	}
	return nil
}

type TransferFrom struct {
	Sender common.Address
	From   common.Address
	To     common.Address
	Value  uint64
}

func (this *TransferFrom) Serialize(w io.Writer) error {
	if err := utils.WriteAddress(w, this.Sender); err != nil {
		return fmt.Errorf("[TransferFrom] serialize sender error:%v", err)
	}
	if err := utils.WriteAddress(w, this.From); err != nil {
		return fmt.Errorf("[TransferFrom] serialize from error:%v", err)
	}
	if err := utils.WriteAddress(w, this.To); err != nil {
		return fmt.Errorf("[TransferFrom] serialize to error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.Value); err != nil {
		return fmt.Errorf("[TransferFrom] serialize value error:%v", err)
	}
	return nil
}

func (this *TransferFrom) Deserialize(r io.Reader) error {
	var err error
	this.Sender, err = utils.ReadAddress(r)
	if err != nil {
		return fmt.Errorf("[TransferFrom] deserialize sender error:%v", err)
	}

	this.From, err = utils.ReadAddress(r)
	if err != nil {
		return fmt.Errorf("[TransferFrom] deserialize from error:%v", err)
	}

	this.To, err = utils.ReadAddress(r)
	if err != nil {
		return fmt.Errorf("[TransferFrom] deserialize to error:%v", err)
	}

	this.Value, err = utils.ReadVarUint(r)
	if err != nil {
		return err
	}
	return nil
}
