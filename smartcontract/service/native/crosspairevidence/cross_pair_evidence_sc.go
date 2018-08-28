package crosspairevidence

import (
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/errors"
	cstates "github.com/mixbee/mixbee/core/states"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"bytes"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/common"
	"strings"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/core/signature"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/smartcontract/service/native/mbc"
	"math"
)

const (
	PUSH_EVIDENCE = "pushEvidence"
	GET_EVIDENCE  = "getCrossPairTx"
)

func InitCrossPairEvidence() {
	native.Contracts[utils.CrossChainPairEvidenceContractAddress] = RegisterCrossPairEvidenceContract
}

func RegisterCrossPairEvidenceContract(native *native.NativeService) {
	native.Register(PUSH_EVIDENCE, PUSH_PAIR_EVIDENCE)
	native.Register(GET_EVIDENCE, GET_CROSS_PAIR_TX)
}

type CTXEntry struct {
	From              string
	To                string
	FromValue         uint64
	ToValue           uint64
	TxHash            string
	ANetWorkId        uint32
	BNetWorkId        uint32
	State             uint32
	SeqId             string
	Type              uint32 //跨链资产类型
	Sig               []byte //验证节点对结果的签名
	Pubk              string //验证节点公钥
	TimeStamp         uint32 //过期时间
	Nonce             uint32 //交易双方的nonce值,必须一样
	VerifyCheckCount  uint32
	ConfrimCheckCount uint32
	ReleaseTxHash     string
}

type CTXPairEntry struct {
	First  *CTXEntry
	Second *CTXEntry
}

func PUSH_PAIR_EVIDENCE(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		log.Warnf("pushEvidence failed: argument %s error %s ", arg0, err.Error())
		return utils.BYTE_FALSE, errors.NewErr("pushEvidence failed: argument 0 error, " + err.Error())
	}

	infoStrs := strings.Split(arg0, ";")
	if len(infoStrs) == 0 {
		log.Warnf("pushEvidence failed: argument %s error %s ", arg0, err.Error())
		return utils.BYTE_FALSE, errors.NewErr("pushEvidence failed: argument invalid")
	}

	infoMap := make(map[string][]byte)
	for _, value := range infoStrs {
		keyValue := strings.Split(value, ":")
		if len(keyValue) != 2 {
			log.Warnf("pushEvidence failed: argument %s error %s ", value, err.Error())
			return utils.BYTE_FALSE, errors.NewErr("pushEvidence failed: argument invalid")
		}
		seqId := keyValue[0]
		vaBuf, err := hex.DecodeString(keyValue[1])
		if err != nil {
			log.Warnf("pushEvidence failed: argument %s error %s ", keyValue[1], err.Error())
			return utils.BYTE_FALSE, errors.NewErr("pushEvidence failed: argument invalid")
		}
		//参数反序列化
		pair := &CTXPairEntry{}
		err = json.Unmarshal(vaBuf, pair)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewErr(fmt.Sprintf("seqId %s json unmarshal error.", seqId))
		}
		//检查签名
		err = checkPairEvidenceSig(pair.First)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewErr(fmt.Sprintf("seqId %s checkPairEvidenceSig. error %s", seqId,err.Error()))
		}
		err = checkPairEvidenceSig(pair.Second)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewErr(fmt.Sprintf("seqId %s checkPairEvidenceSig. error %s", seqId,err.Error()))
		}
		//奖励矿工
		err = awardVerify(native,pair.First)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewErr(fmt.Sprintf("seqId %s awardVerify. error %s", seqId,err.Error()))
		}
		err = awardVerify(native,pair.Second)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewErr(fmt.Sprintf("seqId %s awardVerify. error %s", seqId,err.Error()))
		}
		infoMap[seqId] = vaBuf
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for k, v := range infoMap {
		native.CloneCache.Add(scommon.ST_STORAGE, GenCrossPairEvidenceKey(contract, []byte(k)), &cstates.StorageItem{Value: []byte(v)})
	}
	return utils.BYTE_TRUE, nil
}

func awardVerify(native *native.NativeService,pair *CTXEntry) error {
	pbkByte, err := hex.DecodeString(pair.Pubk)
	if err != nil {
		return errors.NewErr("[awardVerify] pbk invalid,hex decode err. " + err.Error())
	}
	publicKey, err := keypair.DeserializePublicKey(pbkByte)
	if err != nil {
		return errors.NewErr("[awardVerify] pbk invalid,deserialize err. " + err.Error())
	}
	address := types.AddressFromPubKey(publicKey)
	mbgState := &mbc.State{From: utils.MbcContractAddress, To: address, Value: uint64(math.Pow10(9))}
	_, _, err = mbc.TransferForCrossChainContract(native, utils.MbgContractAddress, mbgState)
	if err != nil {
		return errors.NewErr("[awardVerify] ,appCallTransferMbg err. " + err.Error())
	}
	return nil
}

func checkPairEvidenceSig(pair *CTXEntry) error {
	//校验签名信息是否是指定的验证节点
	vpk, err := hex.DecodeString(pair.Pubk)
	if err != nil {
		return fmt.Errorf("VerifyPublicKey hex decode err! err=%s", err)
	}
	publicKey, err := keypair.DeserializePublicKey(vpk)
	if err != nil {
		return fmt.Errorf("VerifyPublicKey DeserializePublicKey err! err=%s", err)
	}
	err = signature.Verify(publicKey, []byte(pair.SeqId), pair.Sig)
	if err != nil {
		return fmt.Errorf("Verify sig info err! err:%s", err.Error())
	}
	return nil
}

func GenCrossPairEvidenceKey(contract common.Address, key []byte) []byte {
	return append(contract[:], key...)
}

func GET_CROSS_PAIR_TX(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	log.Debugf("mixTest args:%s", args)

	arg0, err := serialization.ReadVarBytes(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("add key failed: argument 0 error, " + err.Error())
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	state, err := utils.GetStorageItem(native, GenCrossPairEvidenceKey(contract, arg0))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GET_CROSS_PAIR_TX] get from key error!")
	}
	if state == nil {
		return []byte("not exist"), nil
	}
	return state.Value, nil
}

func appCallTransferMbg(native *native.NativeService, from common.Address, to common.Address, amount uint64) error {
	err := appCallTransfer(native, utils.MbgContractAddress, from, to, amount)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "appCallTransferMbg, appCallTransfer error!")
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
