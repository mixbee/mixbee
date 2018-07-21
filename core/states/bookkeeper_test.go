
package states

import (
	"testing"

	"bytes"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/stretchr/testify/assert"
)

func TestBookkeeper_Deserialize_Serialize(t *testing.T) {
	_, pubKey1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey2, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey3, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey4, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)

	bk := BookkeeperState{
		StateBase:      StateBase{(byte)(1)},
		CurrBookkeeper: []keypair.PublicKey{pubKey1, pubKey2},
		NextBookkeeper: []keypair.PublicKey{pubKey3, pubKey4},
	}

	buf := bytes.NewBuffer(nil)
	bk.Serialize(buf)
	bs := buf.Bytes()

	var bk2 BookkeeperState
	bk2.Deserialize(buf)
	assert.Equal(t, bk, bk2)

	buf = bytes.NewBuffer(bs[:len(bs)-1])
	err := bk2.Deserialize(buf)
	assert.NotNil(t, err)
}
