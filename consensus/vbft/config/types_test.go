

package vconfig

import (
	"encoding/hex"
	"testing"

	"github.com/mixbee/mixbee-crypto/keypair"
)

func TestPubkeyID(t *testing.T) {
	bookkeeper := "120202c924ed1a67fd1719020ce599d723d09d48362376836e04b0be72dfe825e24d81"
	pubKey, err := hex.DecodeString(bookkeeper)
	if err != nil {
		t.Errorf("DecodeString failed: %v", err)
	}
	k, err := keypair.DeserializePublicKey(pubKey)
	if err != nil {
		t.Errorf("DeserializePublicKey failed: %v", err)
	}
	nodeID := PubkeyID(k)
	t.Logf("res: %v\n", nodeID)
}

func TestPubkey(t *testing.T) {
	bookkeeper := "1202027df359dff69eea8dd7d807b669dd9635292b1aae97d03ed32cb36ff30fb7e4d9"
	pubKey, err := hex.DecodeString(bookkeeper)
	if err != nil {
		t.Errorf("DecodeString failed: %v", err)
	}
	k, err := keypair.DeserializePublicKey(pubKey)
	if err != nil {
		t.Errorf("DeserializePublicKey failed: %v", err)
	}
	nodeID := PubkeyID(k)
	publickey, err := Pubkey(nodeID)
	if err != nil {
		t.Errorf("Pubkey failed: %v", err)
	}
	t.Logf("res: %v", publickey)
}
