
package states

import (
	"testing"

	"bytes"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/stretchr/testify/assert"
)

func TestValidatorState_Deserialize_Serialize(t *testing.T) {
	_, pubKey, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)

	vs := ValidatorState{
		StateBase: StateBase{(byte)(1)},
		PublicKey: pubKey,
	}

	buf := bytes.NewBuffer(nil)
	vs.Serialize(buf)
	bs := buf.Bytes()

	var vs2 ValidatorState
	vs2.Deserialize(buf)
	assert.Equal(t, vs, vs2)

	buf = bytes.NewBuffer(bs[:len(bs)-1])
	err := vs2.Deserialize(buf)
	assert.NotNil(t, err)
}
