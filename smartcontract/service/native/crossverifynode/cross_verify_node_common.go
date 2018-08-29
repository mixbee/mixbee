package crossverifynode

import (
	"bytes"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/common"
	"io"
	"fmt"
	"github.com/mixbee/mixbee/smartcontract/service/native/mbc"
	"encoding/binary"
	"github.com/mixbee/mixbee/smartcontract/event"
	"github.com/mixbee/mixbee/common/config"
	"encoding/json"
)

const (
	//constants config
	MIN_DEPOSIT_MBC       = 10000
	WITHDRAW_RELEASE_TIME = 60 * 60 * 48

	// function name
	INIT_NAME                = "init"
	MIN_DEPOSIT_MBC_FUNCTION = "minDepostMbc"
	QUERY_VERIFY_NODE_INFO   = "queryVerifyNodeInfo"
	IS_EXSIT_VERIFY_NODE     = "isExsitVerifyNode"
	REGISTER_VERIFY_NODE     = "registerVerifyNode"
	PAID_DEPOSIT             = "paidDeposit"
	APPLY_WITHDRAW_DEPOSIT   = "applyWithdrawDeposit"
	WITHDRAW_DEPOSIT         = "withdrawDeposit"
	FROZE_DEPOSIT            = "frozeDeposit"
	PUNISHING_DEPOSIT        = "punishingDeposit"
	BLACK_VERIFY_NODE        = "blackNode"
	WHITE_VERIFY_NODE        = "whiteNode"

	//key prefix
	VERIFY_NODE_INFO      = "verifyNodeInfo"
	BLACK_VERIFY_NODES    = "blackVerifyNodes"
	PULISHING_DEPOSIT_SUM = "pulishingDepositSum"
)

const (
	//status
	InitStatus      uint64 = iota
	CanVerifyStatus
	WaitReadyStatus
	BlackStatus
)

type CrossVerifyNodeInfo struct {
	Pbk               string
	Address           string
	Deposit           uint64
	WithdrawDeposit   uint64
	WithDrawStartTime uint64
	FrozeDeposit      uint64
	CurrentStatus     uint64
}

func (this *CrossVerifyNodeInfo) Serialize(w io.Writer) error {
	if err := utils.WriteString(w, this.Pbk); err != nil {
		return fmt.Errorf("[State] serialize Pbk error:%v", err)
	}
	if err := utils.WriteString(w, this.Address); err != nil {
		return fmt.Errorf("[State] serialize address error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.Deposit); err != nil {
		return fmt.Errorf("[State] serialize Deposit error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.WithdrawDeposit); err != nil {
		return fmt.Errorf("[State] serialize WithdrawDeposit error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.WithDrawStartTime); err != nil {
		return fmt.Errorf("[State] serialize WithDrawStartTime error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.FrozeDeposit); err != nil {
		return fmt.Errorf("[State] serialize FrozeDeposit error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.CurrentStatus); err != nil {
		return fmt.Errorf("[State] serialize CurrentStatus error:%v", err)
	}
	return nil
}

func (this *CrossVerifyNodeInfo) Deserialize(r io.Reader) error {
	var err error
	this.Pbk, err = utils.ReadString(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize Pbk error:%v", err)
	}
	this.Address, err = utils.ReadString(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize address error:%v", err)
	}
	this.Deposit, err = utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize Deposit error:%v", err)
	}
	this.WithdrawDeposit, err = utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize WithdrawDeposit error:%v", err)
	}
	this.WithDrawStartTime, err = utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize WithDrawStartTime error:%v", err)
	}
	this.FrozeDeposit, err = utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize FrozeDeposit error:%v", err)
	}
	status, err := utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("[State] deserialize CurrentStatus error:%v", err)
	}
	this.CurrentStatus = status

	return nil
}

func getCrossChainVerifyNodeInfoByPublicKey(native *native.NativeService, pbk string) (*CrossVerifyNodeInfo, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	state, err := utils.GetStorageItem(native, GenCrossChainVerifyNodeInfoKey(contract, pbk))
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[getCrossChainVerifyNodeInfoByPublicKey] CrossVerifyNodeInfo error!")
	}
	if state == nil {
		return nil, nil
	}
	info := CrossVerifyNodeInfo{}
	err = info.Deserialize(bytes.NewBuffer(state.Value))
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[getCrossChainVerifyNodeInfoByPublicKey] CrossVerifyNodeInfo error!")
	}
	return &info, nil
}

func putCrossChainVerifyNodeInfo(native *native.NativeService, info *CrossVerifyNodeInfo) error {

	buf := new(bytes.Buffer)
	err := info.Serialize(buf)
	if err != nil {
		return fmt.Errorf("putCrossChainVerifyNodeInfo Serialize err %s", err.Error())
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	utils.PutBytes(native, GenCrossChainVerifyNodeInfoKey(contract, info.Pbk), buf.Bytes())
	return nil
}

func GenCrossChainVerifyNodeInfoKey(contract common.Address, pbk string) []byte {
	return utils.ConcatKey(contract, []byte(VERIFY_NODE_INFO), []byte(pbk))
}

func appCallTransferMbc(native *native.NativeService, from common.Address, to common.Address, amount uint64) error {
	err := appCallTransfer(native, utils.MbcContractAddress, from, to, amount)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "appCallTransferMbc, appCallTransfer error!")
	}
	return nil
}

func appCallTransfer(native *native.NativeService, contract common.Address, from common.Address, to common.Address, amount uint64) error {
	bf := new(bytes.Buffer)
	var sts []*mbc.State
	sts = append(sts, &mbc.State{
		From:  from,
		To:    to,
		Value: amount,
	})
	transfers := &mbc.Transfers{
		States: sts,
	}
	err := transfers.Serialize(bf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "appCallTransfer, transfers.Serialize error!")
	}

	if _, err := native.NativeCall(contract, "transfer", bf.Bytes()); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "appCallTransfer, appCall error!")
	}
	return nil
}

func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BytesToUint64(buf []byte) uint64 {
	return uint64(binary.BigEndian.Uint64(buf))
}

func AddNotifications(native *native.NativeService, method string, info *CrossVerifyNodeInfo) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	infoB, _ := json.Marshal(info)
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: native.ContextRef.CurrentContext().ContractAddress,
			States:          []interface{}{method, string(infoB)},
		})
}
