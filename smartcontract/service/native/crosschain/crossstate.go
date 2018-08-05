package crosschain

import (
	"github.com/mixbee/mixbee/common"
	"io"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"fmt"
)

const (
	LOCK_STATE   = 0
	REBACK_STATE = 1
	END_STATE    = 2
)

type CrossChainState struct {
	From            common.Address
	To              common.Address
	AValue          uint64
	BValue          uint64
	AChainId        uint32
	BChainId        uint32
	Type            uint32
	Timestamp       uint32
	SeqId           string
	Statue          uint32 //0:lock  1:超时,回退给from  2:成功，转给to
	Nonce           uint32 //跨链双方的匹配随机数
	VerifyPublicKey string //主链验证节点公钥
	Sig             string //主链验证节点签名信息
}

type CrossChainStateResult struct {
	From            string `json:"from"`
	To              string `json:"to"`
	AValue          uint64 `json:"aValue"`
	BValue          uint64 `json:"bValue"`
	AChainId        uint32 `json:"achainId"`
	BChainId        uint32 `json:"bchainId"`
	Type            uint32 `json:"type"`
	Timestamp       uint32 `json:"timestamp"` //有效时间
	SeqId           string `json:"seqId"`
	Statue          uint32 `json:"status"`
	Nonce           uint32 `json:"nonce"`           //跨链双方的匹配随机数
	VerifyPublicKey string `json:"verifyPublicKey"` //主链验证节点公钥
	Sig             string `json:"sig"`
}

func state2Result(state *CrossChainState) *CrossChainStateResult {
	result := &CrossChainStateResult{}
	result.From = state.From.ToBase58()
	result.To = state.To.ToBase58()
	result.AValue = state.AValue
	result.BValue = state.BValue
	result.AChainId = state.AChainId
	result.BChainId = state.BChainId
	result.Type = state.Type
	result.Timestamp = state.Timestamp
	result.SeqId = state.SeqId
	result.Statue = state.Statue
	result.Nonce = state.Nonce
	result.VerifyPublicKey = state.VerifyPublicKey
	result.Sig = state.Sig
	return result
}

func (this *CrossChainState) Serialize(w io.Writer) error {
	if err := utils.WriteAddress(w, this.From); err != nil {
		return fmt.Errorf("[State] serialize from error:%v", err)
	}
	if err := utils.WriteAddress(w, this.To); err != nil {
		return fmt.Errorf("[State] serialize to error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.AValue); err != nil {
		return fmt.Errorf("[State] serialize Avalue error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.BValue); err != nil {
		return fmt.Errorf("[State] serialize Bvalue error:%v", err)
	}
	if err := utils.WriteVarUint(w, uint64(this.AChainId)); err != nil {
		return fmt.Errorf("[State] serialize Bvalue error:%v", err)
	}
	if err := utils.WriteVarUint(w, uint64(this.BChainId)); err != nil {
		return fmt.Errorf("[State] serialize Bvalue error:%v", err)
	}
	if err := utils.WriteVarUint(w, uint64(this.Type)); err != nil {
		return fmt.Errorf("[State] serialize type error:%v", err)
	}
	if err := utils.WriteVarUint(w, uint64(this.Timestamp)); err != nil {
		return fmt.Errorf("[State] serialize timestamp error:%v", err)
	}
	if err := utils.WriteString(w, this.SeqId); err != nil {
		return fmt.Errorf("[State] serialize SeqId error:%v", err)
	}
	if err := utils.WriteVarUint(w, uint64(this.Statue)); err != nil {
		return fmt.Errorf("[State] serialize state error:%v", err)
	}
	if err := utils.WriteVarUint(w, uint64(this.Nonce)); err != nil {
		return fmt.Errorf("[State] serialize state error:%v", err)
	}
	if err := utils.WriteString(w, this.VerifyPublicKey); err != nil {
		return fmt.Errorf("[State] serialize verifyPublicKey  error:%v", err)
	}
	if err := utils.WriteString(w, this.Sig); err != nil {
		return fmt.Errorf("[State] serialize verifyPublicKey  error:%v", err)
	}
	return nil
}

func (this *CrossChainState) Deserialize(r io.Reader) error {
	var err error
	this.From, err = utils.ReadAddress(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize from error:%v", err)
	}
	this.To, err = utils.ReadAddress(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize to error:%v", err)
	}
	this.AValue, err = utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize AValue error:%v", err)
	}
	this.BValue, err = utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize BValue error:%v", err)
	}
	this.AChainId, err = utils.ReadVarUint32(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize AChainId error:%v", err)
	}
	this.BChainId, err = utils.ReadVarUint32(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize BChainId error:%v", err)
	}
	this.Type, err = utils.ReadVarUint32(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize Type error:%v", err)
	}
	this.Timestamp, err = utils.ReadVarUint32(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize Timestamp error:%v", err)
	}
	this.SeqId, err = utils.ReadString(r)
	if err != nil {
		return err
	}
	this.Statue, err = utils.ReadVarUint32(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize state error:%v", err)
	}
	this.Nonce, err = utils.ReadVarUint32(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize nonce error:%v", err)
	}
	this.VerifyPublicKey, err = utils.ReadString(r)
	if err != nil {
		return err
	}
	this.Sig, err = utils.ReadString(r)
	if err != nil {
		return err
	}
	return nil
}

type CrossSeqIds struct {
	SeqIds []string
}

func NewCrossSeqIds() *CrossSeqIds {
	seqs := CrossSeqIds{}
	seqs.SeqIds = []string{}
	return &seqs
}

func (this *CrossSeqIds) Deserialize(r io.Reader) error {

	num, err := utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("CrossSeqIds Deserialize err %s", err.Error())
	}
	for i := 0; i < int(num); i++ {
		str, err := utils.ReadString(r)
		if err != nil {
			return fmt.Errorf("CrossSeqIds Deserialize err %s", err.Error())
		}
		this.SeqIds = append(this.SeqIds, str)
	}

	return nil
}

func (this *CrossSeqIds) Serialize(w io.Writer) error {
	if len(this.SeqIds) == 0 {
		return nil
	}

	num := len(this.SeqIds)
	if err := utils.WriteVarUint(w, uint64(num)); err != nil {
		return fmt.Errorf("[CrossSeqIds] serialize num error:%v", err)
	}
	for _, v := range this.SeqIds {
		err := utils.WriteString(w, v)
		if err != nil {
			return fmt.Errorf("Strings2Bytes err %s", err.Error())
		}
	}
	return nil
}
