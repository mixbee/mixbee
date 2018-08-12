package types

import (
	"github.com/mixbee/mixbee/p2pserver/common"
	"bytes"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
	"fmt"
)

type CrossChainTxCompletedPayload struct {
	SeqId             string
	FirstFrom         string
	FirstReleaseHash  string
	SecondFrom        string
	SecondReleaseHash string
	Type              uint32   //1 release  2 completed
}

//Serialize message payload
func (this *CrossChainTxCompletedPayload) Serialization() ([]byte, error) {

	p := bytes.NewBuffer([]byte{})
	err := serialization.WriteString(p, this.SeqId)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write SeqId error. SeqId :%s", this.SeqId))
	}

	err = serialization.WriteString(p, this.FirstFrom)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write FirstFrom error. FirstFrom:%s", this.FirstFrom))
	}

	err = serialization.WriteString(p, this.FirstReleaseHash)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write FirstReleaseHash error. FirstReleaseHash:%s", this.FirstReleaseHash))
	}

	err = serialization.WriteString(p, this.SecondFrom)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write SecondFrom error. SecondFrom:%s", this.SecondFrom))
	}

	err = serialization.WriteString(p, this.SecondReleaseHash)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write SecondReleaseHash error. SecondReleaseHash:%s", this.SecondReleaseHash))
	}

	err = serialization.WriteUint32(p, this.Type)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write Type error. Type:%d", this.Type))
	}

	return p.Bytes(), nil
}

//Deserialize message payload
func (this *CrossChainTxCompletedPayload) Deserialization(p []byte) error {

	buf := bytes.NewBuffer(p)

	seqId, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read seqId key error")
	}
	this.SeqId = seqId

	firstFrom, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read firstFrom key error")
	}
	this.FirstFrom = firstFrom

	firstReleaseHash, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read firstReleaseHash key error")
	}
	this.FirstReleaseHash = firstReleaseHash

	secondFrom, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read secondFrom key error")
	}
	this.SecondFrom = secondFrom

	secondReleaseHash, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read secondReleaseHash key error")
	}
	this.SecondReleaseHash = secondReleaseHash

	ty, err := serialization.ReadUint32(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read ty key error")
	}
	this.Type = ty

	return nil
}

func (this *CrossChainTxCompletedPayload) CmdType() string {
	return common.CROSSCHAIN_TX_COMPLETED_TYPE
}
