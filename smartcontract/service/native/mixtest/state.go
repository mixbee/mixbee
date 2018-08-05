package mixtest

import (
	"github.com/mixbee/mixbee/common"
	"io"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"fmt"
	"github.com/mixbee/mixbee/errors"
)

type State struct {
	From  common.Address
	Key    string
	Value  string
}

func (this *State) Serialize(w io.Writer) error {
	if err := utils.WriteAddress(w, this.From); err != nil {
		return fmt.Errorf("[State] serialize from error:%v", err)
	}
	if err := utils.WriteString(w, this.Key); err != nil {
		return fmt.Errorf("[State] serialize key error:%v", err)
	}
	if err := utils.WriteString(w, this.Value); err != nil {
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
	this.Key, err = utils.ReadString(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize to error:%v", err)
	}

	this.Value, err = utils.ReadString(r)
	if err != nil {
		return err
	}
	return nil
}

type SetKeys struct {
	States []*State
}

func (this *SetKeys) Serialize(w io.Writer) error {
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

func (this *SetKeys) Deserialize(r io.Reader) error {
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


