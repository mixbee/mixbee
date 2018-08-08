package crosschain

import (
	"github.com/mixbee/mixbee/common/constants"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/errors"
	cstates "github.com/mixbee/mixbee/core/states"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/mbc"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"bytes"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"fmt"
	"github.com/mixbee/mixbee/common"
	"encoding/json"
	"github.com/mixbee/mixbee/common/config"
	"time"
	"strings"
	"encoding/hex"
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/core/signature"
)

func InitCrossChain() {
	native.Contracts[utils.CrossChainContractAddress] = RegisterCrossChainContract
}

func RegisterCrossChainContract(native *native.NativeService) {
	native.Register(mbc.INIT_NAME, Init)
	native.Register(mbc.NAME_NAME, CName)
	native.Register(mbc.SYMBOL_NAME, CSymbol)
	native.Register(CROSS_TRANSFER, CTransfer)
	native.Register(CROSS_QUERY, CQueryBySeqId)
	native.Register(CROSS_HISTORY, CQueryHistory)
	native.Register(CROSS_UNLOCK, CUnlock)
	native.Register(CROSS_RELEASE, CrossRelease)
}

func Init(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	log.Infof("cross chain contract address:%s", contract.ToBase58())
	return utils.BYTE_TRUE, nil
}

func CTransfer(native *native.NativeService) ([]byte, error) {
	//参数反序列化
	state := &CrossChainState{}
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChain] CTransfer deserialize error!")
	}

	//校验seqId是否正确
	ok := CheckSeqId(state)
	if !ok {
		return utils.BYTE_FALSE, errors.NewErr("[CrossChain] CTransfer seqId invalid!")
	}
	//检查seqId是否重复
	info, err := getCrossStateBySeqId(native, state.SeqId)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CTransfer failed:" + err.Error())
	}
	if info != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(nil, errors.ErrDuplicatedSeqId, "duplicated seqId")
	}
	if state.Statue != 0 {
		return utils.BYTE_FALSE, errors.NewErr("[CrossChain] CTransfer state invalid!")
	}
	//检查achainId是否和本网络id相等
	if config.DefConfig.P2PNode.NetworkId != state.AChainId {
		return utils.BYTE_FALSE, errors.NewErr("[CrossChain] achainId not equal node netWorkId !")
	}
	//if state.AChainId == state.BChainId {
	//	return utils.BYTE_FALSE, errors.NewErr("[CrossChain] achainId can not equal bchainId!")
	//}
	//资产转发给智能合约抵押
	contract := native.ContextRef.CurrentContext().ContractAddress
	toAddress, err := BuildDepositAddress(contract, state.From)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChain] build deposit address error!")
	}
	ontState := &mbc.State{From: state.From, To: toAddress, Value: state.AValue}
	_, _, err = mbc.Transfer(native, utils.MbcContractAddress, ontState)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChain] CTransfer asset error!")
	}
	//保存跨链信息 key/value  =  seqId/crossState
	native.CloneCache.Add(scommon.ST_STORAGE, GenCrossChainSeqId(contract, state.SeqId), &cstates.StorageItem{Value: native.Input})
	//保存from地址所有seqId
	err = appendSeqId2Hisory(native, contract, state.From, state.SeqId)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChain] appendSeqId2Hisory error!")
	}
	AddNotifications(native, contract, state)
	return utils.BYTE_TRUE, nil
}

func CrossRelease(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CrossRelease failed: argument 0 error, " + err.Error())
	}

	as := strings.Split(arg0, ":")
	if len(as) != 2 {
		return utils.BYTE_FALSE, errors.NewErr("CrossRelease failed: argument invalid " + arg0)
	}
	seqId := as[0]
	sigStr := as[1]
	sig, err := hex.DecodeString(sigStr)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CrossRelease failed: sign hex decode err " + err.Error())
	}
	//查询seqId对应跨链转移
	info, err := getCrossStateBySeqId(native, seqId)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CQueryBySeqId failed:" + err.Error())
	}

	if info == nil {
		return []byte("not exsit"), nil
	}

	if uint32(time.Now().Unix()) > info.Timestamp {
		return utils.BYTE_FALSE, fmt.Errorf("this tx aleady expire!")
	}

	if info.Statue == REBACK_STATE {
		return utils.BYTE_FALSE, fmt.Errorf("his tx already reback to from!")
	}

	if info.Statue == END_STATE {
		return utils.BYTE_FALSE, fmt.Errorf("this tx already tranfer to !")
	}

	//校验签名信息是否是指定的验证节点
	vpk, err := hex.DecodeString(info.VerifyPublicKey)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("VerifyPublicKey hex decode err! err=%s", err)
	}
	publicKey, err := keypair.DeserializePublicKey(vpk)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("VerifyPublicKey DeserializePublicKey err! err=%s", err)
	}
	err = signature.Verify(publicKey, []byte(info.SeqId), sig)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("Verify sig info err! err:%s", err.Error())
	}

	//把托管资产转给from
	contract := native.ContextRef.CurrentContext().ContractAddress
	toAddress := info.To
	fromAddress, err := BuildDepositAddress(contract, info.From)
	ontState := &mbc.State{From: fromAddress, To: toAddress, Value: info.AValue}
	_, _, err = mbc.TransferForCrossChainContract(native, utils.MbcContractAddress, ontState)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChain] CrossRelease asset error!")
	}
	//更新跨链信息
	info.Statue = 2
	info.Sig = sigStr
	bb := []byte{}
	buf := bytes.NewBuffer(bb)
	err = info.Serialize(buf)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChain] info serialize error!")
	}
	native.CloneCache.Add(scommon.ST_STORAGE, GenCrossChainSeqId(contract, info.SeqId), &cstates.StorageItem{Value: buf.Bytes()})

	return utils.BYTE_TRUE, nil
}

func CUnlock(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CUnlock failed: argument 0 error, " + err.Error())
	}
	//查询seqId对应跨链转移
	info, err := getCrossStateBySeqId(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CQueryBySeqId failed:" + err.Error())
	}

	if info == nil {
		return []byte("not exsit"), nil
	}

	if uint32(time.Now().Unix()) < info.Timestamp {
		return utils.BYTE_FALSE, fmt.Errorf("in lock time,please late try again!")
	}

	if info.Statue == REBACK_STATE {
		return utils.BYTE_FALSE, fmt.Errorf("this tx already reback to from!")
	}

	if info.Statue == END_STATE {
		return utils.BYTE_FALSE, fmt.Errorf("this tx already tranfer to !")
	}

	//把托管资产转给from
	contract := native.ContextRef.CurrentContext().ContractAddress
	toAddress := info.From
	fromAddress, err := BuildDepositAddress(contract, info.From)
	ontState := &mbc.State{From: fromAddress, To: toAddress, Value: info.AValue}
	_, _, err = mbc.TransferForCrossChainContract(native, utils.MbcContractAddress, ontState)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChain] CUnlock asset error!")
	}

	return utils.BYTE_TRUE, nil
}

func CName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.MIXT_NAME), nil
}

func CSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.MIXT_SYMBOL), nil
}

func CQueryBySeqId(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CQueryBySeqId failed: argument 0 error, " + err.Error())
	}

	info, err := getCrossStateBySeqId(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CQueryBySeqId failed:" + err.Error())
	}

	if info == nil {
		return []byte("not exsit"), nil
	}

	jsonStr, err := json.Marshal(state2Result(info))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChainQuery] CrossChainState error!")
	}
	return jsonStr, nil
}

func getCrossStateBySeqId(native *native.NativeService, seqId string) (*CrossChainState, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	state, err := utils.GetStorageItem(native, GenCrossChainSeqId(contract, seqId))
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChainQuery] CrossChainState error!")
	}
	if state == nil {
		return nil, nil
	}
	info := CrossChainState{}
	err = info.Deserialize(bytes.NewBuffer(state.Value))
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[CrossChainQuery] CrossChainState error!")
	}
	return &info, nil
}

func CQueryHistory(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("CQueryBySeqId failed: argument 0 error, " + err.Error())
	}

	from, err := common.AddressFromBase58(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CQueryHistory] AddressFromBase58 error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	seqids := NewCrossSeqIds()
	state, err := utils.GetStorageItem(native, GenCrossChainHistoryKey(contract, from))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CQueryHistory] GetStorageItem error!")
	}
	if state != nil {
		err := seqids.Deserialize(bytes.NewBuffer(state.Value))
		if err != nil {
			return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[CQueryHistory] Deserialize error!")
		}
	}
	result := fmt.Sprintf("%+v", Reverse(seqids.SeqIds))
	return []byte(result), nil
}

func appendSeqId2Hisory(native *native.NativeService, contact, from common.Address, seqId string) error {

	seqids := NewCrossSeqIds()
	//查询历史队列
	state, err := utils.GetStorageItem(native, GenCrossChainHistoryKey(contact, from))
	if err != nil {
		return fmt.Errorf("appendSeqId2Hisory query state err %s", err.Error())
	}
	if state != nil {
		err := seqids.Deserialize(bytes.NewBuffer(state.Value))
		if err != nil {
			return fmt.Errorf("appendSeqId2Hisory err %s", err.Error())
		}
	}
	//新seqId入队列
	seqids.SeqIds = append(seqids.SeqIds, seqId)
	buf := new(bytes.Buffer)
	err = seqids.Serialize(buf)
	if err != nil {
		return fmt.Errorf("appendSeqId2Hisory Serialize err %s", err.Error())
	}
	utils.PutBytes(native, GenCrossChainHistoryKey(contact, from), buf.Bytes())

	return nil
}

func GenCrossChainHistoryKey(contract, from common.Address) []byte {
	return append(contract[:], from[:]...)
}

func GenCrossChainSeqId(contract common.Address, seqId string) []byte {
	return append(contract[:], seqId...)
}
