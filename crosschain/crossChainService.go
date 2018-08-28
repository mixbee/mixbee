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
	p2ptypes "github.com/mixbee/mixbee/p2pserver/message/types"
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

	VerifyNodes    *VerifyNodes
	SubNetNodesMgr *SubChainNetNodes

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
	s.SubNetNodesMgr = NewSubChainNetNodes()

	pidActor := NewCrossChainActor(s)
	pid, err := pidActor.Start()
	if err != nil {
		return nil, fmt.Errorf("crosschaintx actor init error %s", err)
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

	s.SubNetNodesMgr = NewSubChainNetNodes()
	s.SubNetNodesMgr.Init()
}

func (s *CTXPoolServer) Start() {
	log.Infof("cross chain || CTXPoolServer start .....")
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
					log.Infof("seqId= %s is finish push tx to txToMatchPair chan ", pair.First.SeqId)
					s.pairTxPending.push(pair.First)
					s.pairTxPending.push(pair.Second)
				} else {
					log.Infof("send cross tx to p2p server...")
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
				if config.DefConfig.Genesis.ConsensusType != config.CONSENSUS_TYPE_SOLO {
					s.P2pPid.Tell(&p2ptypes.CrossChainTxCompletedPayload{
						SeqId:             pair.First.SeqId,
						FirstFrom:         pair.First.From,
						FirstReleaseHash:  pair.First.ReleaseTxHash,
						SecondFrom:        pair.Second.From,
						SecondReleaseHash: pair.Second.ReleaseTxHash,
						Type:              uint32(1),
					})
				}

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

	log.Infof("cross chain || txEndConfirmHandler len=%v", len(pool.TxList))

	for k, value := range pool.TxList {

		seqId := value.First.SeqId
		value.First.ConfrimCheckCount += 1
		firstPath := server.SubNetNodesMgr.GetSubNetNode(value.First.ANetWorkId)
		firstResult, err := GetCrossChainTxInfoBySeqId(firstPath, seqId)
		if err != nil {
			log.Errorf("cross chain||txEndConfirmHandler||GetCrossChainTxInfoBySeqId error %s ", err)
			continue
		}
		value.First.State = firstResult.Statue

		secondPath := server.SubNetNodesMgr.GetSubNetNode(value.Second.ANetWorkId)
		secondResult, err := GetCrossChainTxInfoBySeqId(secondPath, seqId)
		if err != nil {
			log.Errorf("cross chain||txEndConfirmHandler||GetCrossChainTxInfoBySeqId error %s ", err)
			continue
		}
		value.Second.State = secondResult.Statue

		if firstResult.Statue == 0 && firstResult.Timestamp < uint32(time.Now().Unix()) && value.First.ConfrimCheckCount > 3 {
			value.First.ReleaseTxHash = ""
		}

		if secondResult.Statue == 0 && secondResult.Timestamp < uint32(time.Now().Unix()) && value.First.ConfrimCheckCount > 3 {
			value.Second.ReleaseTxHash = ""
		}

		if value.First.ReleaseTxHash == "" || value.Second.ReleaseTxHash == "" {
			log.Warnf("cross chain || seqid=%s releaseHash invalid,reRelease ", value.First.SeqId)
			value.First.ConfrimCheckCount = 0
			server.txToRelease <- value
			delete(pool.TxList, k)
			continue
		}

		if firstResult.Statue != 0 && secondResult.Statue != 0 {
			log.Infof("cross chain || seqid=%s finished... ", value.First.SeqId)
			delete(pool.TxList, k)
			if config.DefConfig.Genesis.ConsensusType != config.CONSENSUS_TYPE_SOLO {
				server.P2pPid.Tell(&p2ptypes.CrossChainTxCompletedPayload{
					SeqId:             value.First.SeqId,
					FirstFrom:         value.First.From,
					FirstReleaseHash:  value.First.ReleaseTxHash,
					SecondFrom:        value.Second.From,
					SecondReleaseHash: value.Second.ReleaseTxHash,
					Type:              uint32(2),
				})
			}

			//写入智能合约存证
			pushCrossTxEvidence2SmartContract(value,server.VerifyerAccount)
		}
	}
}

func releaseLockToken(s *CTXPoolServer) {
	pool := s.pairTxRelease
	if pool == nil || len(pool.TxList) == 0 {
		return
	}
	pool.Lock()
	defer pool.Unlock()

	log.Infof("cross chain || releaseLockToken len=%v", len(pool.TxList))

	for k, value := range pool.TxList {
		firstSeqId := value.First.SeqId
		secondSeqId := value.Second.SeqId

		firstNetId := value.First.ANetWorkId
		firstPath := s.SubNetNodesMgr.GetSubNetNode(firstNetId)

		secondNetId := value.Second.ANetWorkId
		secondPath := s.SubNetNodesMgr.GetSubNetNode(secondNetId)

		if value.First.ReleaseTxHash == "" {
			firstTxHash, err := pushCrossChainResult(s.VerifyerAccount, firstPath, firstSeqId, value.First.Sig)
			if err != nil {
				log.Errorf("cross chain || pushCrossChainResult first err %v", err)
				continue
			}
			value.First.ReleaseTxHash = firstTxHash

			log.Debugf("cross chain || pushCrossChainResult first success addr=%v,seqId=%v,txHash=%v", firstPath, firstSeqId, firstTxHash)
		}

		if value.Second.ReleaseTxHash == "" {
			secondTxHash, err := pushCrossChainResult(s.VerifyerAccount, secondPath, secondSeqId, value.Second.Sig)
			if err != nil {
				log.Errorf("cross chain || pushCrossChainResult second err %v", err)
				continue
			}
			value.Second.ReleaseTxHash = secondTxHash
			log.Debugf("cross chain || pushCrossChainResult second success addr=%v,seqId=%v,txHash=%v", secondPath, secondSeqId, secondTxHash)
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
	log.Infof("cross chain || txMatchPair check len=%v", len(pool.TxList))

	for k, v := range pool.TxList {
		log.Infof("cross chain || pair first = %#v  second = %#v", v.First, v.Second)
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
		log.Debugf("cross chain || txVerify empty..")
		return
	}

	pool.Lock()
	defer pool.Unlock()

	log.Infof("cross chain || txVerify || check pair cross chain tx poolLen=%v", len(pool.TxList))
	txMap := pool.TxList

	for k, v := range txMap {

		log.Infof("cross chain || txVerify k = %s", k)

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

func (s *CTXPoolServer) CrossTxCompletedHandler(info *p2ptypes.CrossChainTxCompletedPayload) {

	log.Infof("CrossTxCompletedHandler param %v", info)

	seqId := info.SeqId
	//delete txpool
	s.txPool.delete(seqId)

	s.pairTxPending.delete(seqId)

	s.pairTxRelease.delete(seqId)

	if info.Type == 2 {
		s.pairTxEndConfirm.delete(seqId)
	}
}
