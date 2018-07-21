
package signature

import (
	"testing"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/account"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	acc := account.NewAccount("")
	data := []byte{1, 2, 3}
	sig, err := Sign(acc, data)
	assert.Nil(t, err)

	err = Verify(acc.PublicKey, data, sig)
	assert.Nil(t, err)
}

func TestVerifyMultiSignature(t *testing.T) {
	data := []byte{1, 2, 3}
	accs := make([]*account.Account, 0)
	pubkeys := make([]keypair.PublicKey, 0)
	N := 4
	for i := 0; i < N; i++ {
		accs = append(accs, account.NewAccount(""))
	}
	sigs := make([][]byte, 0)

	for _, acc := range accs {
		sig, _ := Sign(acc, data)
		sigs = append(sigs, sig)
		pubkeys = append(pubkeys, acc.PublicKey)
	}

	err := VerifyMultiSignature(data, pubkeys, N, sigs)
	assert.Nil(t, err)

	pubkeys[0], pubkeys[1] = pubkeys[1], pubkeys[0]
	err = VerifyMultiSignature(data, pubkeys, N, sigs)
	assert.Nil(t, err)

}
