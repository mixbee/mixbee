package crossverifynode

import (
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/errors"
	"strconv"
	"bytes"
	"github.com/mixbee/mixbee/common/serialization"
	"encoding/json"
	"github.com/mixbee/mixbee-crypto/keypair"
	"encoding/hex"
	"github.com/mixbee/mixbee/core/types"
	"fmt"
	"strings"
	"github.com/mixbee/mixbee/common"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"math/big"
	ntypes "github.com/mixbee/mixbee/vm/neovm/types"

)

func InitCrossChainVerifyNode() {
	native.Contracts[utils.CrossChainVerifynodeContractAddress] = RegisterCrossChainContract
}

func RegisterCrossChainContract(native *native.NativeService) {
	native.Register(INIT_NAME, Init)
	native.Register(MIN_DEPOSIT_MBC_FUNCTION, MinDepositMbc)
	native.Register(QUERY_VERIFY_NODE_INFO, QueryVerifyNodeInfo)
	native.Register(IS_EXSIT_VERIFY_NODE, IsExsitVerifyNode)
	native.Register(REGISTER_VERIFY_NODE, RegisterVerifyNode)
	native.Register(PAID_DEPOSIT, PaidDeposit)
	native.Register(APPLY_WITHDRAW_DEPOSIT, ApplyWithdrawDeposit)
	native.Register(WITHDRAW_DEPOSIT, WithdrawDeposit)
	native.Register(FROZE_DEPOSIT, FrozeDeposit)
	native.Register(PUNISHING_DEPOSIT, PunishingDeposit)
	native.Register(BLACK_VERIFY_NODE, BlackNode)
	native.Register(WHITE_VERIFY_NODE, WhiteNode)
}

func Init(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	// check if init is already execute
	pulishDepositBytes, err := native.CloneCache.Get(scommon.ST_STORAGE, utils.ConcatKey(contract, []byte(PULISHING_DEPOSIT_SUM)))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getGovernanceView, get governanceViewBytes error!")
	}
	if pulishDepositBytes != nil {
		return utils.BYTE_FALSE, errors.NewErr("init is already executed!")
	}
	utils.PutBytes(native, utils.ConcatKey(contract, []byte(PULISHING_DEPOSIT_SUM)), Uint64ToBytes(0))
	log.Infof("cross chain verify node contract address:%s", contract.ToBase58())
	return utils.BYTE_TRUE, nil
}

func MinDepositMbc(native *native.NativeService) ([]byte, error) {
	return ntypes.BigIntToBytes(big.NewInt(int64(MIN_DEPOSIT_MBC))), nil
}

func IsExsitVerifyNode(native *native.NativeService) ([]byte, error) {
	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("isExsitVerifyNode failed: argument 0 error, " + err.Error())
	}
	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("isExsitVerifyNode failed:" + err.Error())
	}

	if info == nil {
		return []byte("false"), nil
	}

	return []byte("true"), nil
}

func QueryVerifyNodeInfo(native *native.NativeService) ([]byte, error) {
	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("QueryVerifyNodeInfo failed: argument 0 error, " + err.Error())
	}
	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("QueryVerifyNodeInfo failed:" + err.Error())
	}

	if info == nil {
		return nil, nil
	}

	jsonStr, err := json.Marshal(info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[QueryVerifyNodeInfo] CrossChainVerifyNodeInfo error!")
	}
	return jsonStr, nil
}

func RegisterVerifyNode(native *native.NativeService) ([]byte, error) {

	//参数反序列化
	info := &CrossVerifyNodeInfo{}
	if err := info.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[RegisterVerifyNode] CrossVerifyNodeInfo deserialize error!")
	}
	if info.Pbk == "" {
		return utils.BYTE_FALSE, errors.NewErr("[RegisterVerifyNode] pbk invalid!")
	}
	pbkByte, err := hex.DecodeString(info.Pbk)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("[RegisterVerifyNode] pbk invalid,hex decode err. " + err.Error())
	}
	publicKey, err := keypair.DeserializePublicKey(pbkByte)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("[RegisterVerifyNode] pbk invalid,deserialize err. " + err.Error())
	}
	address := types.AddressFromPubKey(publicKey)
	info.Address = address.ToBase58()

	//check witness
	err = utils.ValidateOwner(native, address)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "validateOwner, checkWitness error!")
	}

	storeInfo, err := getCrossChainVerifyNodeInfoByPublicKey(native, info.Pbk)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "validateExsit, check is exsit error!")
	}
	if storeInfo != nil {
		return utils.BYTE_FALSE, fmt.Errorf("[RegisterVerifyNode] pbk %s is exsit ", info.Pbk)
	}

	info.CurrentStatus = InitStatus
	info.Deposit = 0
	info.FrozeDeposit = 0
	info.WithdrawDeposit = 0
	info.WithDrawStartTime = 0
	err = putCrossChainVerifyNodeInfo(native, info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "putCrossChainVerifyNodeInfo error!")
	}

	return utils.BYTE_TRUE, nil
}

func PaidDeposit(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("PaidDeposit failed: argument 0 error, " + err.Error())
	}

	as := strings.Split(arg0, ":")
	if len(as) != 2 {
		return utils.BYTE_FALSE, errors.NewErr("PaidDeposit failed: argument invalid " + arg0)
	}
	pbk := as[0]
	amountStr := as[1]
	deposit, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "deposit invliad!")
	}
	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, pbk)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getCrossChainVerifyNodeInfoByPublicKey error!")
	}
	if info == nil {
		return utils.BYTE_FALSE, errors.NewErr("not found verifyNode pbk=" + pbk)
	}

	if info.CurrentStatus == BlackStatus {
		return utils.BYTE_FALSE, fmt.Errorf("pbk=%s verifyNode is in blackNode list", arg0)
	}

	err = appCallTransferMbc(native, native.Tx.Payer, utils.CrossChainVerifynodeContractAddress,deposit)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "deposit transfer error!")
	}

	info.Deposit = info.Deposit + deposit

	if info.Deposit >= MIN_DEPOSIT_MBC {
		info.CurrentStatus = CanVerifyStatus
	} else {
		info.CurrentStatus = WaitReadyStatus
	}

	err = putCrossChainVerifyNodeInfo(native, info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "putCrossChainVerifyNodeInfo error!")
	}

	return utils.BYTE_TRUE, nil
}

func ApplyWithdrawDeposit(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("ApplyWithdrawDeposit failed: argument 0 error, " + err.Error())
	}

	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getCrossChainVerifyNodeInfoByPublicKey error!")
	}
	if info == nil {
		return utils.BYTE_FALSE, errors.NewErr("not found verifyNode pbk=" + arg0)
	}
	if info.CurrentStatus == BlackStatus {
		return utils.BYTE_FALSE, fmt.Errorf("pbk=%s verifyNode is in blackNode list", arg0)
	}

	info.WithdrawDeposit = info.Deposit
	info.Deposit = 0
	info.WithDrawStartTime = uint64(native.Time)
	info.CurrentStatus = WaitReadyStatus

	err = putCrossChainVerifyNodeInfo(native, info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "putCrossChainVerifyNodeInfo error!")
	}

	return utils.BYTE_TRUE, nil
}

func WithdrawDeposit(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("WithdrawDeposit failed: argument 0 error, " + err.Error())
	}

	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getCrossChainVerifyNodeInfoByPublicKey error!")
	}
	if info == nil {
		return utils.BYTE_FALSE, errors.NewErr("not found verifyNode pbk=" + arg0)
	}
	if info.CurrentStatus == BlackStatus {
		return utils.BYTE_FALSE, fmt.Errorf("pbk=%s verifyNode is in blackNode list", arg0)
	}

	if info.WithdrawDeposit == 0 {
		return utils.BYTE_FALSE, fmt.Errorf("pbk=%s verifyNode withdraw deposit is zero", arg0)
	}

	if info.WithDrawStartTime+WITHDRAW_RELEASE_TIME >= uint64(native.Time) {
		return utils.BYTE_FALSE, fmt.Errorf("pbk=%s verifyNode withdraw release time is %d", info.Pbk, info.WithDrawStartTime+WITHDRAW_RELEASE_TIME)
	}

	address, err := common.AddressFromBase58(info.Address)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, fmt.Sprintf("WithdrawDeposit||address=%s to object errr", info.Address))
	}

	err = appCallTransferMbc(native, utils.CrossChainVerifynodeContractAddress, address, info.WithdrawDeposit)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "WithdrawDeposit transfer error!")
	}

	err = putCrossChainVerifyNodeInfo(native, info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "putCrossChainVerifyNodeInfo error!")
	}

	return utils.BYTE_TRUE, nil
}

func FrozeDeposit(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("PaidDeposit failed: argument 0 error, " + err.Error())
	}

	as := strings.Split(arg0, ":")
	if len(as) != 2 {
		return utils.BYTE_FALSE, errors.NewErr("PaidDeposit failed: argument invalid " + arg0)
	}
	pbk := as[0]
	amountStr := as[1]
	deposit, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "frozeDeposit invliad!")
	}

	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, pbk)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getCrossChainVerifyNodeInfoByPublicKey error!")
	}
	if info == nil {
		return utils.BYTE_FALSE, errors.NewErr("not found verifyNode pbk=" + pbk)
	}

	if deposit > info.Deposit {
		return utils.BYTE_FALSE, fmt.Errorf("frozeDeposit %v great than deposit %v", deposit, info.Deposit)
	}

	info.FrozeDeposit = deposit
	info.Deposit = info.Deposit - deposit
	if info.Deposit < MIN_DEPOSIT_MBC {
		info.CurrentStatus = WaitReadyStatus
	}

	err = putCrossChainVerifyNodeInfo(native, info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "putCrossChainVerifyNodeInfo error!")
	}

	return utils.BYTE_TRUE, nil
}

func PunishingDeposit(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("PunishingDeposit failed: argument 0 error, " + err.Error())
	}

	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getCrossChainVerifyNodeInfoByPublicKey error!")
	}
	if info == nil {
		return utils.BYTE_FALSE, errors.NewErr("not found verifyNode pbk=" + arg0)
	}

	if info.FrozeDeposit == 0 {
		return utils.BYTE_FALSE, fmt.Errorf("pbk=%s verifyNode withdraw deposit is zero", arg0)
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	pulishDepositBytes, err := utils.GetStorageItem(native, utils.ConcatKey(contract, []byte(PULISHING_DEPOSIT_SUM)))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getGovernanceView, get governanceViewBytes error!")
	}
	if pulishDepositBytes != nil {
		return utils.BYTE_FALSE, errors.NewErr("init is already executed!")
	}

	pulishDeposit := BytesToUint64(pulishDepositBytes.Value)
	utils.PutBytes(native, utils.ConcatKey(contract, []byte(PULISHING_DEPOSIT_SUM)), Uint64ToBytes(pulishDeposit))
	pulishDeposit = pulishDeposit + info.FrozeDeposit
	info.FrozeDeposit = 0
	err = putCrossChainVerifyNodeInfo(native, info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "putCrossChainVerifyNodeInfo error!")
	}

	return utils.BYTE_TRUE, nil
}

func BlackNode(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("BlackNode failed: argument 0 error, " + err.Error())
	}

	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getCrossChainVerifyNodeInfoByPublicKey error!")
	}
	if info == nil {
		return utils.BYTE_FALSE, errors.NewErr("not found verifyNode pbk=" + arg0)
	}

	info.FrozeDeposit = info.Deposit + info.WithdrawDeposit
	info.Deposit = 0
	info.WithdrawDeposit = 0
	info.CurrentStatus = BlackStatus
	err = putCrossChainVerifyNodeInfo(native, info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "putCrossChainVerifyNodeInfo error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	utils.PutBytes(native, utils.ConcatKey(contract, []byte(BLACK_VERIFY_NODES)), []byte(info.Pbk))

	return utils.BYTE_TRUE, nil
}

func WhiteNode(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("WhiteNode failed: argument 0 error, " + err.Error())
	}

	info, err := getCrossChainVerifyNodeInfoByPublicKey(native, arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "getCrossChainVerifyNodeInfoByPublicKey error!")
	}
	if info == nil {
		return utils.BYTE_FALSE, errors.NewErr("not found verifyNode pbk=" + arg0)
	}

	info.Deposit = info.FrozeDeposit
	if info.Deposit >= MIN_DEPOSIT_MBC {
		info.CurrentStatus = CanVerifyStatus
	} else {
		info.CurrentStatus = WaitReadyStatus
	}
	err = putCrossChainVerifyNodeInfo(native, info)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "putCrossChainVerifyNodeInfo error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CloneCache.Delete(scommon.ST_STORAGE, utils.ConcatKey(contract, []byte(BLACK_VERIFY_NODES), []byte(info.Pbk)))
	return utils.BYTE_TRUE, nil
}
