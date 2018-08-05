package crosschain

import (
	"sync"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee-crypto/keypair"
)


type CTXEntry struct {
	From  			string
	To    			string
	FromValue 		uint64
	ToValue 		uint64
	TxHash 			string
	ANetWorkId 		uint32
	BNetWorkId 		uint32
	State           uint32  //0 待确认  1 打包成功
	SeqId           string
	Type            uint  //跨链资产类型
	Sig             []byte //验证节点对结果的签名
	Pubk            keypair.PublicKey //验证节点公钥
	TimeStamp       uint32  //过期时间
	Nonce           uint32  //交易双方的nonce值,必须一样
}

type CTXMatchPool struct {
	sync.RWMutex
	TxList map[string] []*CTXEntry // Transactions which have been verified
}

type CTXPairEntrys  []*CTXPairEntry

type CTXPairEntry struct {
	First *CTXEntry
	Second *CTXEntry
}

// Init creates a new transaction pool to gather.
func (tp *CTXMatchPool) Init() {
	tp.Lock()
	defer tp.Unlock()
	tp.TxList = make(map[string] []*CTXEntry)
}

func (tp *CTXMatchPool) push(entry *CTXEntry) {
	tp.Lock()
	defer tp.Unlock()

	if _,ok := tp.TxList[entry.SeqId];!ok {
		tp.TxList[entry.SeqId] = []*CTXEntry{}
	}

	list := tp.TxList[entry.SeqId]
	list = append(list,entry)
	tp.TxList[entry.SeqId] = list
	log.Infof("CTXMatchPool push success. len = %v",len(tp.TxList[entry.SeqId]))
}