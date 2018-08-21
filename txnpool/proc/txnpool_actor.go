

package proc

import (
	"reflect"

	"github.com/mixbee/mixbee/mixbee-eventbus/actor"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/core/ledger"
	tx "github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/events/message"
	hComm "github.com/mixbee/mixbee/http/base/common"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/smartcontract/service/neovm"
	tc "github.com/mixbee/mixbee/txnpool/common"
	"github.com/mixbee/mixbee/validator/types"
)

// NewTxActor creates an actor to handle the transaction-based messages from
// network and http
func NewTxActor(s *TXPoolServer) *TxActor {
	a := &TxActor{}
	a.setServer(s)
	return a
}

// NewTxPoolActor creates an actor to handle the messages from the consensus
func NewTxPoolActor(s *TXPoolServer) *TxPoolActor {
	a := &TxPoolActor{}
	a.setServer(s)
	return a
}

// NewVerifyRspActor creates an actor to handle the verified result from validators
func NewVerifyRspActor(s *TXPoolServer) *VerifyRspActor {
	a := &VerifyRspActor{}
	a.setServer(s)
	return a
}

// isBalanceEnough checks if the tranactor has enough to cover gas cost
func isBalanceEnough(address common.Address, gas uint64) bool {
	balance, err := hComm.GetContractBalance(0, utils.MbgContractAddress, address)
	if err != nil {
		log.Debugf("failed to get contract balance %s err %v",
			address.ToHexString(), err)
		return false
	}
	return balance >= gas
}

// TxnActor: Handle the low priority msg from P2P and API
type TxActor struct {
	server *TXPoolServer
}

// handleTransaction handles a transaction from network and http
func (ta *TxActor) handleTransaction(sender tc.SenderType, self *actor.PID,
	txn *tx.Transaction) {
	ta.server.increaseStats(tc.RcvStats)
	if len(txn.ToArray()) > tc.MAX_TX_SIZE {
		log.Debugf("handleTransaction: reject a transaction due to size over 1M")
		return
	}

	if ta.server.getTransaction(txn.Hash()) != nil {
		log.Debugf("handleTransaction: transaction %x already in the txn pool",
			txn.Hash())

		ta.server.increaseStats(tc.DuplicateStats)
	} else if ta.server.getTransactionCount() >= tc.MAX_CAPACITY {
		log.Infof("handleTransaction: transaction pool is full for tx %x",
			txn.Hash())

		ta.server.increaseStats(tc.FailureStats)
	} else {
		if _, overflow := common.SafeMul(txn.GasLimit, txn.GasPrice); overflow {
			log.Debugf("handleTransaction: gasLimit %v, gasPrice %v overflow",
				txn.GasLimit, txn.GasPrice)
			return
		}

		if !txn.SystemTx {
			if txn.GasLimit < config.DefConfig.Common.GasLimit ||
				txn.GasPrice < ta.server.getGasPrice() {
				log.Debugf("handleTransaction: invalid gasLimit %v, gasPrice %v",
					txn.GasLimit, txn.GasPrice)
				return
			}

			if txn.TxType == tx.Deploy && txn.GasLimit < neovm.CONTRACT_CREATE_GAS {
				log.Debugf("handleTransaction: deploy tx invalid gasLimit %v, gasPrice %v",
					txn.GasLimit, txn.GasPrice)
				return
			}
		}

		if ta.server.preExec && sender != tc.HttpSender {
			result, err := ledger.DefLedger.PreExecuteContract(txn)
			if err != nil {
				log.Debugf("handleTransaction: failed to preExecuteContract tx %x err %v",
					txn.Hash(), err)
			}
			if txn.GasLimit < result.Gas {
				log.Debugf("handleTransaction: transaction's gasLimit %d is less than preExec gasLimit %d",
					txn.GasLimit, result.Gas)
				return
			}
			gas, overflow := common.SafeMul(txn.GasPrice, result.Gas)
			if overflow {
				log.Debugf("handleTransaction: gasPrice %d preExec gasLimit %d overflow",
					txn.GasPrice, result.Gas)
				return
			}
			if !isBalanceEnough(txn.Payer, gas) {
				log.Debugf("handleTransaction: transactor %s has no balance enough to cover gas cost %d",
					txn.Payer.ToHexString(), gas)
				return
			}
			log.Debugf("handleTransaction: tx %x preExec success", txn.Hash())
		}
		<-ta.server.slots
		ta.server.assignTxToWorker(txn, sender)
	}
}

// Receive implements the actor interface
func (ta *TxActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("txpool-tx actor started and be ready to receive tx msg")

	case *actor.Stopping:
		log.Warn("txpool-tx actor stopping")

	case *actor.Restarting:
		log.Warn("txpool-tx actor restarting")

	case *tc.TxReq:
		sender := msg.Sender

		log.Debugf("txpool-tx actor receives tx from %v ", sender.Sender())

		ta.handleTransaction(sender, context.Self(), msg.Tx)

	case *tc.GetTxnReq:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives getting tx req from %v", sender)

		res := ta.server.getTransaction(msg.Hash)
		if sender != nil {
			sender.Request(&tc.GetTxnRsp{Txn: res},
				context.Self())
		}

	case *tc.GetTxnStats:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives getting tx stats from %v", sender)

		res := ta.server.getStats()
		if sender != nil {
			sender.Request(&tc.GetTxnStatsRsp{Count: res},
				context.Self())
		}

	case *tc.CheckTxnReq:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives checking tx req from %v", sender)

		res := ta.server.checkTx(msg.Hash)
		if sender != nil {
			sender.Request(&tc.CheckTxnRsp{Ok: res},
				context.Self())
		}

	case *tc.GetTxnStatusReq:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives getting tx status req from %v", sender)

		res := ta.server.getTxStatusReq(msg.Hash)
		if sender != nil {
			if res == nil {
				sender.Request(&tc.GetTxnStatusRsp{Hash: msg.Hash,
					TxStatus: nil}, context.Self())
			} else {
				sender.Request(&tc.GetTxnStatusRsp{Hash: res.Hash,
					TxStatus: res.Attrs}, context.Self())
			}
		}

	case *tc.GetTxnCountReq:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives getting tx count req from %v", sender)

		res := ta.server.getTxCount()
		if sender != nil {
			sender.Request(&tc.GetTxnCountRsp{Count: res},
				context.Self())
		}

	default:
		log.Debugf("txpool-tx actor: unknown msg %v type %v", msg, reflect.TypeOf(msg))
	}
}

func (ta *TxActor) setServer(s *TXPoolServer) {
	ta.server = s
}

// TxnPoolActor: Handle the high priority request from Consensus
type TxPoolActor struct {
	server *TXPoolServer
}

// Receive implements the actor interface
func (tpa *TxPoolActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("txpool actor started and be ready to receive txPool msg")

	case *actor.Stopping:
		log.Warn("txpool actor stopping")

	case *actor.Restarting:
		log.Warn("txpool actor Restarting")

	case *tc.GetTxnPoolReq:
		sender := context.Sender()

		log.Debugf("txpool actor receives getting tx pool req from %v", sender)

		res := tpa.server.getTxPool(msg.ByCount, msg.Height)
		if sender != nil {
			sender.Request(&tc.GetTxnPoolRsp{TxnPool: res}, context.Self())
		}

	case *tc.GetPendingTxnReq:
		sender := context.Sender()

		log.Debugf("txpool actor receives getting pedning tx req from %v", sender)

		res := tpa.server.getPendingTxs(msg.ByCount)
		if sender != nil {
			sender.Request(&tc.GetPendingTxnRsp{Txs: res}, context.Self())
		}

	case *tc.VerifyBlockReq:
		sender := context.Sender()

		log.Debugf("txpool actor receives verifying block req from %v", sender)

		tpa.server.verifyBlock(msg, sender)

	case *message.SaveBlockCompleteMsg:
		sender := context.Sender()

		log.Debugf("txpool actor receives block complete event from %v", sender)

		if msg.Block != nil {
			tpa.server.cleanTransactionList(msg.Block.Transactions, msg.Block.Header.Height)
		}

	default:
		log.Debugf("txpool actor: unknown msg %v type %v", msg, reflect.TypeOf(msg))
	}
}

func (tpa *TxPoolActor) setServer(s *TXPoolServer) {
	tpa.server = s
}

// VerifyRspActor: Handle the response from the validators
type VerifyRspActor struct {
	server *TXPoolServer
}

// Receive implements the actor interface
func (vpa *VerifyRspActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("txpool-verify actor: started and be ready to receive validator's msg")

	case *actor.Stopping:
		log.Warn("txpool-verify actor: stopping")

	case *actor.Restarting:
		log.Warn("txpool-verify actor: Restarting")

	case *types.RegisterValidator:
		log.Debugf("txpool-verify actor:: validator %v connected", msg.Sender)
		vpa.server.registerValidator(msg)

	case *types.UnRegisterValidator:
		log.Debugf("txpool-verify actor:: validator %d:%v disconnected", msg.Type, msg.Id)

		vpa.server.unRegisterValidator(msg.Type, msg.Id)

	case *types.CheckResponse:
		log.Debug("txpool-verify actor:: Receives verify rsp message")

		vpa.server.assignRspToWorker(msg)

	default:
		log.Debugf("txpool-verify actor:Unknown msg %v type %v", msg, reflect.TypeOf(msg))
	}
}

func (vpa *VerifyRspActor) setServer(s *TXPoolServer) {
	vpa.server = s
}
