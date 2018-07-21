

package consensus

import (
	"github.com/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/consensus/dbft"
	"github.com/mixbee/mixbee/consensus/solo"
	"github.com/mixbee/mixbee/consensus/vbft"
)

type ConsensusService interface {
	Start() error
	Halt() error
	GetPID() *actor.PID
}

const (
	CONSENSUS_DBFT = "dbft"
	CONSENSUS_SOLO = "solo"
	CONSENSUS_VBFT = "vbft"
)

func NewConsensusService(consensusType string, account *account.Account, txpool *actor.PID, ledger *actor.PID, p2p *actor.PID) (ConsensusService, error) {
	if consensusType == "" {
		consensusType = CONSENSUS_DBFT
	}
	var consensus ConsensusService
	var err error
	switch consensusType {
	case CONSENSUS_DBFT:
		consensus, err = dbft.NewDbftService(account, txpool, p2p)
	case CONSENSUS_SOLO:
		consensus, err = solo.NewSoloService(account, txpool)
	case CONSENSUS_VBFT:
		consensus, err = vbft.NewVbftServer(account, txpool, p2p)
	}
	log.Infof("ConsensusType:%s", consensusType)
	return consensus, err
}
