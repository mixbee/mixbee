package crosschain

import (
	"sync"
	"github.com/mixbee/mixbee/common/log"
)

type CTXEntry struct {
	From          string
	To            string
	FromValue     uint64
	ToValue       uint64
	TxHash        string
	ANetWorkId    uint32
	BNetWorkId    uint32
	State         uint32 //0 待确认  1 打包成功
	SeqId         string
	Type          uint32 //跨链资产类型
	Sig           []byte //验证节点对结果的签名
	Pubk          string //验证节点公钥
	TimeStamp     uint32 //过期时间
	Nonce         uint32 //交易双方的nonce值,必须一样
	CheckCount    uint32
	ReleaseTxHash string
}

type CTXMatchPool struct {
	sync.RWMutex
	TxList map[string]*CTXPairEntry // Transactions which have been verified
}

type CTXPairEntry struct {
	First  *CTXEntry
	Second *CTXEntry
}

// Init creates a new transaction pool to gather.
func (tp *CTXMatchPool) Init() {
	tp.Lock()
	defer tp.Unlock()
	tp.TxList = make(map[string]*CTXPairEntry)
}

func (tp *CTXMatchPool) delete (seqId string) {
	tp.Lock()
	defer tp.Unlock()
	delete(tp.TxList,seqId)
}

func (tp *CTXMatchPool) push(entry *CTXEntry) {
	tp.Lock()
	defer tp.Unlock()
	log.Infof("cross chain CTXMatchPool push seqId=%s from=%s",entry.SeqId,entry.From)
	if _, ok := tp.TxList[entry.SeqId]; !ok {
		tp.TxList[entry.SeqId] = &CTXPairEntry{}
	}

	pair := tp.TxList[entry.SeqId]
	//check repeat tx
	if pair.First != nil && pair.First.From == entry.From {
		log.Infof("repeat cross tx seqId=%s,from=%s",entry.SeqId,entry.From)
		return
	}
	if pair.Second != nil && pair.Second.From == entry.From {
		log.Infof("repeat cross tx seqId=%s,from=%s",entry.SeqId,entry.From)
		return
	}

	if pair.First == nil {
		pair.First = entry
	} else if pair.Second == nil {
		pair.Second = entry
	} else {
		log.Errorf("CTXMatchPool push error")
	}
	tp.TxList[entry.SeqId] = pair
}
