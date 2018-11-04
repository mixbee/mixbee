

package util

import (
	"crypto/sha256"
	"errors"
	"io"

	"github.com/mixbee/mixbee-crypto/keypair"
	s "github.com/mixbee/mixbee-crypto/signature"
	"github.com/mixbee/mixbee/common/log"
	mixErrors "github.com/mixbee/mixbee/errors"
	"golang.org/x/crypto/ripemd160"
)

type ECDsaCrypto struct {
}

func (c *ECDsaCrypto) Hash160(message []byte) []byte {
	temp := sha256.Sum256(message)
	md := ripemd160.New()
	io.WriteString(md, string(temp[:]))
	hash := md.Sum(nil)
	return hash
}

func (c *ECDsaCrypto) Hash256(message []byte) []byte {
	temp := sha256.Sum256(message)
	f := sha256.Sum256(temp[:])
	return f[:]
}

func (c *ECDsaCrypto) VerifySignature(message []byte, signature []byte, pubkey []byte) (bool, error) {

	log.Debugf("message: %x", message)
	log.Debugf("signature: %x", signature)
	log.Debugf("pubkey: %x", pubkey)

	pk, err := keypair.DeserializePublicKey(pubkey)
	if err != nil {
		return false, mixErrors.NewDetailErr(errors.New("[ECDsaCrypto], deserializing public key failed."), mixErrors.ErrNoCode, "")
	}

	sig, err := s.Deserialize(signature)
	ok := s.Verify(pk, message, sig)
	if !ok {
		return false, mixErrors.NewDetailErr(errors.New("[ECDsaCrypto], VerifySignature failed."), mixErrors.ErrNoCode, "")
	}

	return true, nil
}
