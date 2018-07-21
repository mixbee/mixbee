

package vote

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/core/genesis"
	"github.com/mixbee/mixbee/core/types"
)

func GetValidators(txs []*types.Transaction) ([]keypair.PublicKey, error) {
	// TODO implement vote
	return genesis.GenesisBookkeepers, nil
}
