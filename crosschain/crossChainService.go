package crosschain

import (
	"sync"
	"time"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/config"
	"fmt"
	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/core/signature"
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
	httpactor "github.com/mixbee/mixbee/http/base/actor"
	"github.com/mixbee/mixbee-crypto/keypair"
	"encoding/hex"
)

var CtxServer *CTXPoolServer

// TXPoolServer contains all api to external modules
type CTXPoolServer struct {
	mu sync.RWMutex   // Sync mutex
	wg sync.WaitGroup // Worker sync
	//workers          []ctxPoolWorker // Worker pool
	txPool        *CTXMatchPool // The tx pool that holds the valid transaction
	pairTxPending *CTXMatchPool
	txToMatchPair chan *CTXPairEntry
	//slots            chan struct{} // The limited slots for the new transaction
	pairTxRelease    *CTXMatchPool
	txToRelease      chan *CTXPairEntry
	pairTxEndConfirm *CTXMatchPool
	txToEndConfirm   chan *CTXPairEntry
	VerifyerAccount  *account.Account

	VerifyNodes *VerifyNodes

	P2pPid *actor.PID
	Pid    *actor.PID
}

func NewCTxPoolServer(num uint8, acc *account.Account, p2pPid *actor.PID) (*actor.PID, error) {
	s := &CTXPoolServer{}
	s.init(num, p2pPid)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.VerifyerAccount = acc
	s.VerifyNodes = NewVerifyNodes()

	pidActor := NewCrossChainActor(s)
	pid, err := pidActor.Start()
	if err != nil {
		return nil, fmt.Errorf("crosschain actor init error %s", err)
	}
	s.Pid = pid

	CtxServer = s
	return pid, nil
}

// init initializes the server with the configured settings
func (s *CTXPoolServer) init(num uint8, p2pPid *actor.PID) {
	// Initial txnPool
	s.txPool = &CTXMatchPool{}
	s.txPool.Init()
	s.P2pPid = p2pPid
	s.pairTxPending = &CTXMatchPool{}
	s.pairTxPending.Init()
	s.pairTxRelease = &CTXMatchPool{}
	s.pairTxRelease.Init()
	s.pairTxEndConfirm = &CTXMatchPool{}
	s.pairTxEndConfirm.Init()

	s.txToMatchPair = make(chan *CTXPairEntry, 5000)
	s.txToRelease = make(chan *CTXPairEntry, 5000)
	s.txToEndConfirm = make(chan *CTXPairEntry, 5000)
}

func (s *CTXPoolServer) Start() {
	log.Infof("CTXPoolServer start .....")
	ticker := time.NewTicker(config.DEFAULT_CROSS_CHAIN_VERIFY_TIME * time.Second)
	for {
		select {
		case <-ticker.C:
			go txVerify(s)
			go txMatchPair(s)
			go releaseLockToken(s)
			go txEndConfirmHandler(s)
		case pair, ok := <-s.txToMatchPair:
			if ok {
				if pair.First != nil && pair.Second != nil {
					s.pairTxPending.push(pair.First)
					s.pairTxPending.push(pair.Second)
				} else {
					pushSigedCrossTx2OtherNode(pair, s)
				}
			}
		case pair, ok := <-s.txToRelease:
			if ok {
				s.pairTxRelease.push(pair.First)
				s.pairTxRelease.push(pair.Second)
			}
		case pair, ok := <-s.txToEndConfirm:
			if ok {
				s.pairTxEndConfirm.push(pair.First)
				s.pairTxEndConfirm.push(pair.Second)
			}
		}
	}
}
func txEndConfirmHandler(server *CTXPoolServer) {

	pool := server.pairTxEndConfirm
	if pool == nil || len(pool.TxList) == 0 {
		return
	}

	pool.Lock()
	defer pool.Unlock()

	log.Infof("cross chain txEndConfirmHandler len=%v", len(pool.TxList))

	subChainNode := config.DefConfig.CrossChain.SubChainNode

	for k, value := range pool.TxList {
		firstPath := subChainNode[value.First.ANetWorkId]
		firstHash := value.First.ReleaseTxHash
		firstState, err := GetTxStateByHash(firstPath[0], firstHash)
		if err != nil {
			log.Errorf("GetTxStateByHash addr=%v,hash=%v err=%v", firstPath[0], firstHash, err)
			continue
		}

		secondPath := subChainNode[value.First.ANetWorkId]
		secondHash := value.Second.ReleaseTxHash
		secondState, err := GetTxStateByHash(secondPath[0], secondHash)
		if err != nil {
			log.Errorf("GetTxStateByHash addr=%v,hash=%v err=%v", secondPath[0], secondHash, err)
			continue
		}

		if firstState == 1 && secondState == 1 {
			delete(pool.TxList, k)
			continue
		}

		if firstState != 1 {
			value.First.ReleaseTxHash = ""
		}

		if secondState != 1 {
			value.Second.ReleaseTxHash = ""
		}

		server.txToRelease <- value
		delete(pool.TxList, k)
	}
}

func releaseLockToken(s *CTXPoolServer) {
	pool := s.pairTxRelease
	if pool == nil || len(pool.TxList) == 0 {
		return
	}
	pool.Lock()
	defer pool.Unlock()

	log.Infof("cross chain releaseLockToken len=%v", len(pool.TxList))

	subChainNode := config.DefConfig.CrossChain.SubChainNode

	for k, value := range pool.TxList {
		firstSeqId := value.First.SeqId
		secondSeqId := value.Second.SeqId

		firstNetId := value.First.ANetWorkId
		firstPath := subChainNode[firstNetId][0]

		secondNetId := value.Second.ANetWorkId
		secondPath := subChainNode[secondNetId][0]

		if value.First.ReleaseTxHash == "" {
			firstTxHash, err := pushCrossChainResult(s.VerifyerAccount, firstPath, firstSeqId, value.First.Sig)
			if err != nil {
				log.Errorf("pushCrossChainResult first err %v", err)
				continue
			}
			value.First.ReleaseTxHash = firstTxHash

			log.Debugf("pushCrossChainResult first success addr=%v,seqId=%v,txHash=%v", firstPath, firstSeqId, firstTxHash)
		}

		if value.Second.ReleaseTxHash == "" {
			secondTxHash, err := pushCrossChainResult(s.VerifyerAccount, secondPath, secondSeqId, value.Second.Sig)
			if err != nil {
				log.Errorf("pushCrossChainResult second err %v", err)
				continue
			}
			value.Second.ReleaseTxHash = secondTxHash
			log.Debugf("pushCrossChainResult second success addr=%v,seqId=%v,txHash=%v", secondPath, secondSeqId, secondTxHash)
		}

		delete(pool.TxList, k)
		s.txToEndConfirm <- value
	}
}

//校验交易双方是否打包
func txMatchPair(s *CTXPoolServer) {

	pool := s.pairTxPending
	if pool == nil || len(pool.TxList) == 0 {
		return
	}
	pool.Lock()
	defer pool.Unlock()
	log.Infof("txMatchPair check len=%v", len(pool.TxList))

	for k, v := range pool.TxList {
		if v.First != nil && v.Second != nil {
			s.txToRelease <- v
			delete(pool.TxList, k)
		}
	}
}

//把交易匹配打包
func txVerify(s *CTXPoolServer) {

	pool := s.txPool
	if pool == nil || len(pool.TxList) == 0 {
		log.Debugf("txVerify empty..")
		return
	}

	pool.Lock()
	defer pool.Unlock()

	log.Infof("txVerify || check pair cross chain tx poolLen=%v", len(pool.TxList))
	txMap := pool.TxList

	for k, v := range txMap {

		log.Infof("txVerify k = %s", k)

		first := v.First
		if first != nil {
			ok, expire := checkCrossChainTxBySeqId(first)
			if expire {
				delete(txMap, k)
				pool.TxList = txMap
				continue
			}

			if !ok {
				continue
			}
			//sign cross chain tx
			sigDate := first.SeqId
			sig, err := signature.Sign(s.VerifyerAccount, []byte(sigDate))
			if err != nil {
				log.Errorf("txMatchPair signature err %v", err.Error())
				continue
			}
			first.Sig = sig
		}

		second := v.Second
		if second != nil {
			ok, expire := checkCrossChainTxBySeqId(second)
			if expire {
				delete(txMap, k)
				pool.TxList = txMap
				return
			}

			if !ok {
				return
			}
			//sign cross chain tx

			sigDate := second.SeqId
			sig, err := signature.Sign(s.VerifyerAccount, []byte(sigDate))
			if err != nil {
				log.Errorf("txMatchPair signature err %v", err.Error())
				continue
			}
			second.Sig = sig
		}

		s.txToMatchPair <- v
		delete(txMap, k)
		pool.TxList = txMap
	}
}

// init initializes the server with the configured settings
func (s *CTXPoolServer) PushCtxToPool(rsq *httpactor.PushCrossChainTxRsq) error {
	log.Infof("cross chain PushCtxToPool params=%#v", rsq)

	entry := &CTXEntry{
		From:       rsq.From,
		To:         rsq.To,
		FromValue:  rsq.FromValue,
		ToValue:    rsq.ToValue,
		ANetWorkId: rsq.ANetWorkId,
		BNetWorkId: rsq.BNetWorkId,
		TxHash:     rsq.TxHash,
		SeqId:      rsq.SeqId,
		TimeStamp:  rsq.TimeStamp,
		Nonce:      rsq.Nonce,
		Pubk:       rsq.Pubk,
	}

	if s.IsVerifyNode(rsq.Pubk) {
		s.txPool.push(entry)
	} else {
		log.Warnf("cross chain tx push err node. nodePublicKey=%s tx=%#v", s.VerifyerAccount.PublicKey, rsq)
	}

	return nil
}

func (s *CTXPoolServer) IsVerifyNode(pbk string) bool {
	if s.VerifyerAccount == nil {
		return false
	}

	bb := keypair.SerializePublicKey(s.VerifyerAccount.PublicKey)
	publicKey := hex.EncodeToString(bb)
	return pbk == publicKey
}
