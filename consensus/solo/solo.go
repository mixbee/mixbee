

package solo

import (
	"fmt"
	"reflect"
	"time"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	actorTypes "github.com/mixbee/mixbee/consensus/actor"
	"github.com/mixbee/mixbee/core/ledger"
	"github.com/mixbee/mixbee/core/signature"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/events"
	"github.com/mixbee/mixbee/events/message"
	"github.com/mixbee/mixbee/validator/increment"
)

/*
*Simple consensus for solo node in test environment.
 */
const ContextVersion uint32 = 0

type SoloService struct {
	Account          *account.Account
	poolActor        *actorTypes.TxPoolActor
	incrValidator    *increment.IncrementValidator
	existCh          chan interface{}
	genBlockInterval time.Duration
	pid              *actor.PID
	sub              *events.ActorSubscriber
}

func NewSoloService(bkAccount *account.Account, txpool *actor.PID) (*SoloService, error) {
	service := &SoloService{
		Account:          bkAccount,
		poolActor:        &actorTypes.TxPoolActor{Pool: txpool},
		incrValidator:    increment.NewIncrementValidator(10),
		genBlockInterval: time.Duration(config.DefConfig.Genesis.SOLO.GenBlockTime) * time.Second,
	}

	props := actor.FromProducer(func() actor.Actor {
		return service
	})

	pid, err := actor.SpawnNamed(props, "consensus_solo")
	service.pid = pid
	service.sub = events.NewActorSubscriber(pid)

	return service, err
}

func (self *SoloService) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Restarting:
		log.Info("solo actor restarting")
	case *actor.Stopping:
		log.Info("solo actor stopping")
	case *actor.Stopped:
		log.Info("solo actor stopped")
	case *actor.Started:
		log.Info("solo actor started")
	case *actor.Restart:
		log.Info("solo actor restart")
	case *actorTypes.StartConsensus:
		if self.existCh != nil {
			log.Info("consensus have started")
			return
		}

		self.sub.Subscribe(message.TOPIC_SAVE_BLOCK_COMPLETE)

		timer := time.NewTicker(self.genBlockInterval)
		self.existCh = make(chan interface{})
		go func() {
			defer timer.Stop()
			existCh := self.existCh
			for {
				select {
				case <-timer.C:
					self.pid.Tell(&actorTypes.TimeOut{})
				case <-existCh:
					return
				}
			}
		}()
	case *actorTypes.StopConsensus:
		if self.existCh != nil {
			close(self.existCh)
			self.existCh = nil
			self.incrValidator.Clean()
			self.sub.Unsubscribe(message.TOPIC_SAVE_BLOCK_COMPLETE)
		}
	case *message.SaveBlockCompleteMsg:
		log.Infof("solo actor receives block complete event. block height=%d txnum=%d", msg.Block.Header.Height, len(msg.Block.Transactions))
		self.incrValidator.AddBlock(msg.Block)

	case *actorTypes.TimeOut:
		err := self.genBlock()
		if err != nil {
			log.Errorf("Solo genBlock error %s", err)
		}
	default:
		log.Info("solo actor: Unknown msg ", msg, "type", reflect.TypeOf(msg))
	}
}

func (self *SoloService) GetPID() *actor.PID {
	return self.pid
}

func (self *SoloService) Start() error {
	self.pid.Tell(&actorTypes.StartConsensus{})
	return nil
}

func (self *SoloService) Halt() error {
	self.pid.Tell(&actorTypes.StopConsensus{})
	return nil
}

func (self *SoloService) genBlock() error {
	block, err := self.makeBlock()
	if err != nil {
		return fmt.Errorf("makeBlock error %s", err)
	}

	err = ledger.DefLedger.AddBlock(block)
	if err != nil {
		return fmt.Errorf("genBlock DefLedgerPid.RequestFuture Height:%d error:%s", block.Header.Height, err)
	}
	return nil
}

func (self *SoloService) makeBlock() (*types.Block, error) {
	log.Debug()
	owner := self.Account.PublicKey
	nextBookkeeper, err := types.AddressFromBookkeepers([]keypair.PublicKey{owner})
	if err != nil {
		return nil, fmt.Errorf("GetBookkeeperAddress error:%s", err)
	}
	prevHash := ledger.DefLedger.GetCurrentBlockHash()
	height := ledger.DefLedger.GetCurrentBlockHeight()

	validHeight := height

	start, end := self.incrValidator.BlockRange()

	if height+1 == end {
		validHeight = start
	} else {
		self.incrValidator.Clean()
		log.Infof("increment validator block height %v != ledger block height %v", int(end)-1, height)
	}

	log.Infof("current block height %v, increment validator block cache range: [%d, %d)", height, start, end)

	txs := self.poolActor.GetTxnPool(true, validHeight)

	transactions := make([]*types.Transaction, 0, len(txs))
	for _, txEntry := range txs {
		// TODO optimize to use height in txentry
		if err := self.incrValidator.Verify(txEntry.Tx, validHeight); err == nil {
			transactions = append(transactions, txEntry.Tx)
		}
	}

	txHash := []common.Uint256{}
	for _, t := range transactions {
		txHash = append(txHash, t.Hash())
	}
	txRoot := common.ComputeMerkleRoot(txHash)

	blockRoot := ledger.DefLedger.GetBlockRootWithNewTxRoot(txRoot)
	header := &types.Header{
		Version:          ContextVersion,
		PrevBlockHash:    prevHash,
		TransactionsRoot: txRoot,
		BlockRoot:        blockRoot,
		Timestamp:        uint32(time.Now().Unix()),
		Height:           height + 1,
		ConsensusData:    common.GetNonce(),
		NextBookkeeper:   nextBookkeeper,
	}
	block := &types.Block{
		Header:       header,
		Transactions: transactions,
	}

	blockHash := block.Hash()

	sig, err := signature.Sign(self.Account, blockHash[:])
	if err != nil {
		return nil, fmt.Errorf("[Signature],Sign error:%s.", err)
	}

	block.Header.Bookkeepers = []keypair.PublicKey{owner}
	block.Header.SigData = [][]byte{sig}
	return block, nil
}
