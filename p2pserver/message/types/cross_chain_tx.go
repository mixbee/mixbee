package types

import (
	"github.com/mixbee/mixbee/p2pserver/common"
	"bytes"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
	"fmt"
)

type CrossChainTxInfoPayload struct {
	From       string
	To         string
	FromValue  uint64
	ToValue    uint64
	TxHash     string
	ANetWorkId uint32
	BNetWorkId uint32
	SeqId      string
	Type       uint32              //跨链资产类型
	Sig        []byte            //验证节点对结果的签名
	Pubk       string //验证节点公钥
	TimeStamp  uint32            //过期时间
	Nonce      uint32            //交易双方的nonce值,必须一样
}


//Serialize message payload
func (this *CrossChainTxInfoPayload) Serialization() ([]byte, error) {

	p := bytes.NewBuffer([]byte{})
	err := serialization.WriteString(p, this.From)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write from error. from :%s", this.From))
	}

	err = serialization.WriteString(p, this.To)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write to error. to:%s", this.To))
	}

	err = serialization.WriteUint64(p, this.FromValue)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write fromValue error. fromValue:%d", this.FromValue))
	}

	err = serialization.WriteUint64(p, this.ToValue)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write toValue error. toValue:%d", this.ToValue))
	}

	err = serialization.WriteString(p, this.TxHash)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write txhash error. txhash:%s", this.TxHash))
	}

	err = serialization.WriteUint32(p, this.ANetWorkId)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write ANetWorkId error. ANetWorkId:%d", this.ANetWorkId))
	}

	err = serialization.WriteUint32(p, this.BNetWorkId)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write BNetWorkId error. BNetWorkId:%d", this.BNetWorkId))
	}

	err = serialization.WriteString(p, this.SeqId)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write SeqId error. SeqId:%s", this.SeqId))
	}

	err = serialization.WriteUint32(p, this.Type)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write Type error. Type:%d", this.Type))
	}

	err = serialization.WriteVarBytes(p, this.Sig)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write Sig error. Sig:%x", this.Sig))
	}

	err = serialization.WriteString(p, this.Pubk)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write Pubk error. Pubk:%s", this.Pubk))
	}

	err = serialization.WriteUint32(p, this.TimeStamp)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write TimeStamp error. TimeStamp:%d", this.TimeStamp))
	}

	err = serialization.WriteUint32(p, this.Nonce)
	if err != nil {
		return nil,errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write Nonce error. Nonce:%d", this.Nonce))
	}

	return p.Bytes(),nil
}

//Deserialize message payload
func (this *CrossChainTxInfoPayload) Deserialization(p []byte) error {

	buf := bytes.NewBuffer(p)

	from, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read from key error")
	}
	this.From = from

	to, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read to key error")
	}
	this.To = to

	fromValue, err := serialization.ReadUint64(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read fromValue key error")
	}
	this.FromValue = fromValue

	toValue, err := serialization.ReadUint64(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read toValue key error")
	}
	this.ToValue = toValue

	txHash, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read txHash key error")
	}
	this.TxHash = txHash

	anid, err := serialization.ReadUint32(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read anid key error")
	}
	this.ANetWorkId = anid

	bnid, err := serialization.ReadUint32(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read bnid key error")
	}
	this.BNetWorkId = bnid

	seqId, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read seqId key error")
	}
	this.SeqId = seqId

	ty, err := serialization.ReadUint32(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read type key error")
	}
	this.Type = ty

	sig,err := serialization.ReadVarBytes(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read sig key error")
	}
	this.Sig = sig

	pbk, err := serialization.ReadString(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read pbk key error")
	}
	this.Pubk = pbk

	time, err := serialization.ReadUint32(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read timestamp key error")
	}
	this.TimeStamp = time

	nonce, err := serialization.ReadUint32(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read nonce key error")
	}
	this.Nonce = nonce

	return nil
}

func (this *CrossChainTxInfoPayload) CmdType() string {
	return common.CROSSCHAIN_TX_TYPE
}
