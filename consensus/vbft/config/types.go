

package vconfig

import (
	"encoding/hex"
	"fmt"

	"github.com/mixbee/mixbee-crypto/keypair"
)

// PubkeyID returns a marshaled representation of the given public key.
func PubkeyID(pub keypair.PublicKey) string {
	nodeid := hex.EncodeToString(keypair.SerializePublicKey(pub))
	return nodeid
}

func Pubkey(nodeid string) (keypair.PublicKey, error) {
	pubKey, err := hex.DecodeString(nodeid)
	if err != nil {
		return nil, err
	}
	pk, err := keypair.DeserializePublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("deserialize failed: %s", err)
	}
	return pk, err
}
