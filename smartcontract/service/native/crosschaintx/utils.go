package crosschaintx

import (
	"encoding/binary"
	"crypto/sha256"
	"encoding/hex"
	"github.com/mixbee/mixbee/common"
	"golang.org/x/crypto/ripemd160"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/smartcontract/event"
	"encoding/json"
	"sort"
	"fmt"
)

const (
	CROSS_TRANSFER = "crossTranfer"
	CROSS_QUERY    = "crossQuery"
	CROSS_HISTORY  = "crossHistory"
	CROSS_UNLOCK  = "crossUnlock"  //托管资产回退到from
	CROSS_RELEASE  = "crossRelease"  //托管资产转给to
)

/**
构建跨链抵押地址
由跨链智能合约地址和抵押人地址计算得到
contract + from  --> hash256  -->  hash160  --> address
 */
func BuildDepositAddress(contract, from common.Address) (common.Address,error ) {
	accountIdByte := sha256.Sum256(append(contract[:],from[:]...))
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(accountIdByte[:])
	hashBytes := ripemd160Hash.Sum(nil)
	addr,err := common.AddressParseFromBytes(hashBytes)
	return addr,err
}

func CheckSeqId(state *CrossChainState) bool {
	seqId := GetSeqId(state)
	return  seqId == state.SeqId
}

func GetSeqId(state *CrossChainState) string {

	sList := []string{}
	sList = append(sList,state.From.ToHexString())
	sList = append(sList,state.To.ToHexString())
	sList = append(sList,string(Int64ToBytes(state.AValue)))
	sList = append(sList,string(Int64ToBytes(state.BValue)))
	sList = append(sList,string(Int32ToBytes(state.AChainId)))
	sList = append(sList,string(Int32ToBytes(state.BChainId)))
	sList = append(sList,string(Int32ToBytes(state.Type)))
	sList = append(sList,string(Int32ToBytes(state.Nonce)))
	//sList = append(sList,string(Int32ToBytes(state.Timestamp)))

	sort.Sort(sort.StringSlice(sList))
	var list []byte
	for _,v := range sList {
		list = append(list,v...)
	}
	bvar := sha256.Sum256(list)
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(bvar[:])
	hashBytes := ripemd160Hash.Sum(nil)
	fmt.Println("GetSeqId=",hex.EncodeToString(hashBytes[:]))
	return hex.EncodeToString(hashBytes[:])
}

func Int64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func Int32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func Reverse(str []string) []string {
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	return str
}

func AddNotifications(native *native.NativeService, contract common.Address, state *CrossChainState) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	jsonStr,_ := json.Marshal(state2Result(state))
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{CROSS_TRANSFER,string(jsonStr)},
		})
}