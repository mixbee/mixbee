

package stateful

import (
	"github.com/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/core/ledger"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/validator/db"
	vatypes "github.com/mixbee/mixbee/validator/types"
	"reflect"
)

// Validator is an interface for tx validation actor
type Validator interface {
	Register(poolId *actor.PID)
	UnRegister(poolId *actor.PID)
	VerifyType() vatypes.VerifyType
}

type validator struct {
	pid       *actor.PID
	id        string
	bestBlock db.BestBlock
}

// NewValidator returns Validator for stateful check of tx
func NewValidator(id string) (Validator, error) {

	validator := &validator{id: id}
	props := actor.FromProducer(func() actor.Actor {
		return validator
	})

	pid, err := actor.SpawnNamed(props, id)
	validator.pid = pid
	return validator, err
}

func (self *validator) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("stateful-validator: started and be ready to receive txn")
	case *actor.Stopping:
		log.Info("stateful-validator: stopping")
	case *actor.Restarting:
		log.Info("stateful-validator: restarting")
	case *vatypes.CheckTx:
		log.Debugf("stateful-validator: receive tx %x", msg.Tx.Hash())
		sender := context.Sender()
		height := ledger.DefLedger.GetCurrentBlockHeight()

		errCode := errors.ErrNoError
		hash := msg.Tx.Hash()

		exist, err := ledger.DefLedger.IsContainTransaction(hash)
		if err != nil {
			log.Warn("query db error:", err)
			errCode = errors.ErrUnknown
		} else if exist {
			errCode = errors.ErrDuplicatedTx
		}

		response := &vatypes.CheckResponse{
			WorkerId: msg.WorkerId,
			Type:     self.VerifyType(),
			Hash:     msg.Tx.Hash(),
			Height:   height,
			ErrCode:  errCode,
		}

		sender.Tell(response)
	case *vatypes.UnRegisterAck:
		context.Self().Stop()
	case *types.Block:

		//bestBlock, _ := self.db.GetBestBlock()
		//if bestBlock.Height+1 < msg.Header.Height {
		//	// add sync block request
		//} else if bestBlock.Height+1 == msg.Header.Height {
		//	self.db.PersistBlock(msg)
		//}

	default:
		log.Info("stateful-validator: unknown msg ", msg, "type", reflect.TypeOf(msg))
	}

}

func (self *validator) VerifyType() vatypes.VerifyType {
	return vatypes.Stateful
}

func (self *validator) Register(poolId *actor.PID) {
	poolId.Tell(&vatypes.RegisterValidator{
		Sender: self.pid,
		Type:   self.VerifyType(),
		Id:     self.id,
	})
}

func (self *validator) UnRegister(poolId *actor.PID) {
	poolId.Tell(&vatypes.UnRegisterValidator{
		Id: self.id,
	})
}
