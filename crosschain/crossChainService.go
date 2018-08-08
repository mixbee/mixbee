package crosschain

import (
	"sync"
	"time"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/cmd/utils"
	"encoding/json"
	"fmt"
	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/core/signature"
)

var CtxServer *CTXPoolServer

// TXPoolServer contains all api to external modules
type CTXPoolServer struct {
	mu              sync.RWMutex // Sync mutex
	muPool          sync.RWMutex
	muPending       sync.RWMutex
	muVerifyed      sync.RWMutex
	wg              sync.WaitGroup  // Worker sync
	workers         []ctxPoolWorker // Worker pool
	txPool          *CTXMatchPool   // The tx pool that holds the valid transaction
	pairTxPending   CTXPairEntrys
	slots           chan struct{} // The limited slots for the new transaction
	txToPending     chan *CTXPairEntry
	pairTxVerifyed  CTXPairEntrys
	txToverify      chan *CTXPairEntry
	VerifyerAccount *account.Account
}

func NewCTxPoolServer(num uint8, acc *account.Account) *CTXPoolServer {
	s := &CTXPoolServer{}
	s.init(num)
	s.mu.Lock()
	defer s.mu.Unlock()
	CtxServer = s
	s.VerifyerAccount = acc
	return s
}

// init initializes the server with the configured settings
func (s *CTXPoolServer) init(num uint8) {
	// Initial txnPool
	s.txPool = &CTXMatchPool{}
	s.txPool.Init()

	s.pairTxPending = CTXPairEntrys{}
	s.pairTxVerifyed = CTXPairEntrys{}

	s.txToPending = make(chan *CTXPairEntry, 10000)
	s.txToverify = make(chan *CTXPairEntry, 10000)
}

func (s *CTXPoolServer) Start() {
	log.Infof("CTXPoolServer start .....")
	ticker := time.NewTicker(config.DEFAULT_CROSS_CHAIN_VERIFY_TIME * time.Second)
	for {
		select {
		case <-ticker.C:
			go txToMatchPair(s)
			go verifyBlock(s)
			go releaseLockToken(s)
		case pair, ok := <-s.txToPending:
			if ok {
				s.pairTxPending = append(s.pairTxPending, pair)
			}
		case pair, ok := <-s.txToverify:
			if ok {
				s.pairTxVerifyed = append(s.pairTxVerifyed, pair)
			}
		}
	}
}

func releaseLockToken(s *CTXPoolServer) {
	pool := s.pairTxVerifyed
	if pool == nil || len(pool) == 0 {
		return
	}
	s.muVerifyed.Lock()
	defer s.muVerifyed.Unlock()
	log.Infof("cross chain releaseLockToken len=%v", len(pool))
	var indexs []int
	subChainNode := config.DefConfig.CrossChain.SubChainNode
	for index, value := range pool {
		firstSeqId := value.First.SeqId
		secondSeqId := value.Second.SeqId

		firstNetId := value.First.ANetWorkId
		firstPath := subChainNode[firstNetId][0]

		secondNetId := value.Second.ANetWorkId
		secondPath := subChainNode[secondNetId][0]

		firstTxHash, err := pushCrossChainResult(s.VerifyerAccount, firstPath, firstSeqId, value.First.Sig)
		if err != nil {
			log.Errorf("pushCrossChainResult first err %v", err)
			continue
		}

		log.Infof("pushCrossChainResult first success addr=%v,seqId=%v,txHash=%v", firstPath, firstSeqId, firstTxHash)

		secondTxHash, err := pushCrossChainResult(s.VerifyerAccount, secondPath, secondSeqId, value.Second.Sig)
		if err != nil {
			log.Errorf("pushCrossChainResult second err %v", err)
			continue
		}

		log.Infof("pushCrossChainResult second success addr=%v,seqId=%v,txHash=%v", secondPath, secondSeqId, secondTxHash)

		indexs = append(indexs, index)
	}

	//delete
	log.Infof("releaseLockToken delete release tx.. len=%v", len(indexs))
	for e := range indexs {
		pool = append(pool[:e], pool[e+1:]...)
		s.pairTxVerifyed = pool
	}
}

//校验交易双方是否打包
func verifyBlock(s *CTXPoolServer) {

	pool := s.pairTxPending
	if pool == nil || len(pool) == 0 {
		return
	}
	s.muPending.Lock()
	defer s.muPending.Unlock()
	log.Infof("verifyBlock check len=%v", len(pool))
	var indexs []int
	subChainNode := config.DefConfig.CrossChain.SubChainNode
	for index, v := range pool {
		first := v.First
		second := v.Second

		firstState := v.First.State
		if firstState == 0 {
			firstHash := first.TxHash
			firstPath := subChainNode[first.ANetWorkId]
			firstState, err := GetTxStateByHash(firstPath[0], firstHash)
			if err != nil {
				log.Errorf("GetTxStateByHash addr=%v,hash=%v err=%v", firstPath[0], firstHash, err)
			}
			if firstState == 1 {
				v.First.State = firstState
				sigDate := v.First.SeqId
				sig, err := signature.Sign(s.VerifyerAccount, []byte(sigDate))
				if err != nil {
					log.Errorf("verifyBlock signature err %v", err.Error())
				}
				v.First.Sig = sig
			}
		}

		secondState := v.Second.State
		if secondState == 0 {
			secondHash := second.TxHash
			secondPath := subChainNode[second.ANetWorkId]
			secondState, err := GetTxStateByHash(secondPath[0], secondHash)
			if err != nil {
				log.Errorf("GetTxStateByHash addr=%v,hash=%v err=%v", secondPath[0], secondHash, err)
			}

			if secondState == 1 {
				v.Second.State = secondState
				sigDate := v.Second.SeqId
				sig, err := signature.Sign(s.VerifyerAccount, []byte(sigDate))
				if err != nil {
					log.Errorf("verifyBlock signature err %v", err.Error())
				}
				v.Second.Sig = sig
			}
		}

		if firstState == 1 && secondState == 1 {
			s.txToverify <- v
			indexs = append(indexs, index)
		}
	}

	//delete
	log.Infof("verifyBlock delete match tx.. len=%v", len(indexs))
	for e := range indexs {
		pool = append(pool[:e], pool[e+1:]...)
		s.pairTxPending = pool
	}
}

func pushCrossChainResult(signer *account.Account, addr, seqId string, sig []byte) (string, error) {
	log.Infof("releaseLockToken||pushCrossChainResult addr=%v,seqid=%v", addr, seqId)
	result, err := utils.CrossChainReleaseAssetByMainChain(signer, addr, seqId, sig)
	if err != nil {
		return "", err
	}
	log.Infof("pushCrossChainResult %s", result)
	return result, nil
}

func GetTxStateByHash(addr, hash string) (uint32, error) {
	//log.Infof("verifyBlock||GetTxStateByHash addr=%v,hash=%v",addr,hash)
	result, err := utils.SendRpcRequestWithAddr(addr, "getsmartcodeevent", []interface{}{hash})
	if err != nil {
		return 0, err
	}

	re, err := Json2map(result)
	if err != nil {
		return 0, err
	}
	if re["State"] == nil {
		return 0, nil
	}

	bb, ok := re["State"].(float64)
	if !ok {
		log.Errorf("crossChainService||GetTxStateByHash err %s", result)
		return 0, nil
	}
	return uint32(bb), nil
}

func Json2map(param []byte) (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal(param, &result); err != nil {
		return nil, err
	}
	return result, nil
}

//把交易匹配打包
func txToMatchPair(s *CTXPoolServer) {

	pool := s.txPool
	if pool == nil || len(pool.TxList) == 0 {
		log.Infof("txToMatchPair empty..")
		return
	}

	s.muPool.Lock()
	defer s.muPool.Unlock()

	log.Infof("txToMatchPair || check pair cross chain tx poolLen=%v", len(pool.TxList))
	txMap := pool.TxList

	for k, v := range txMap {
		if len(v) < 2 {
			continue
		}

		log.Infof("txToMatchPair k = %s", k)

		first := v[0]
		second := v[1]

		if first.From != second.To {
			v = append(v[1:])
			v = append(v, first)
			continue
		}

		//v = append(v[1:])
		v = append(v[2:])
		pair := CTXPairEntry{First: first, Second: second}
		s.txToPending <- &pair

		if len(v) == 0 {
			delete(txMap, k)
			pool.TxList = txMap
		}
	}
}

// init initializes the server with the configured settings
func (s *CTXPoolServer) PushCtxToPool(params []interface{}) error {
	log.Infof("cross chain PushCtxToPool params=%#v", params)
	from, ok := params[0].(string)
	if !ok {
		return fmt.Errorf("param-0 invalid %v", params[0])
	}
	to, ok := params[1].(string)
	if !ok {
		return fmt.Errorf("param-1 invalid %v", params[1])
	}
	fromV, ok := params[2].(float64)
	if !ok {
		return fmt.Errorf("param-2 invalid %v", params[2])
	}
	toV, ok := params[3].(float64)
	if !ok {
		return fmt.Errorf("param-3 invalid %v", params[3])
	}
	aNId, ok := params[4].(float64)
	if !ok {
		return fmt.Errorf("param-4 invalid %v", params[4])
	}
	bNId, ok := params[5].(float64)
	if !ok {
		return fmt.Errorf("param-5 invalid %v", params[5])
	}
	txHash, ok := params[6].(string)
	if !ok {
		return fmt.Errorf("param-6 invalid %v", params[6])
	}
	seqId, ok := params[7].(string)
	if !ok {
		return fmt.Errorf("param-7 invalid %v", params[7])
	}
	timestamp, ok := params[8].(float64)
	if !ok {
		return fmt.Errorf("param-8 invalid %v", params[8])
	}
	nonce, ok := params[9].(float64)
	if !ok {
		return fmt.Errorf("param-9 invalid %v", params[9])
	}
	publickKey, ok := params[10].(string)
	if !ok {
		return fmt.Errorf("param-7 invalid %v", params[10])
	}

	nodeinfo, ok := config.DefConfig.CrossChain.SubChainNode[uint32(aNId)]
	if !ok || len(nodeinfo) == 0 {
		return fmt.Errorf("invalid networkId %v,subchain not register in mainchain", aNId)
	}

	nodeinfo, ok = config.DefConfig.CrossChain.SubChainNode[uint32(bNId)]
	if !ok || len(nodeinfo) == 0 {
		return fmt.Errorf("invalid networkId %v,subchain not register in mainchain", bNId)
	}

	entry := &CTXEntry{
		From:       from,
		To:         to,
		FromValue:  uint64(fromV),
		ToValue:    uint64(toV),
		ANetWorkId: uint32(aNId),
		BNetWorkId: uint32(bNId),
		TxHash:     txHash,
		SeqId:      seqId,
		TimeStamp:  uint32(timestamp),
		Nonce:      uint32(nonce),
		Pubk:       publickKey,
	}

	s.txPool.push(entry)

	return nil
}
