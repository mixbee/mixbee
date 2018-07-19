package address

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/account"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

func AddressFromPubKey(pubkey keypair.PublicKey) {
}


func AddressFromVmCode(code []byte) account.Address {
	var addr account.Address
	temp := sha256.Sum256(code)
	md := ripemd160.New()
	md.Write(temp[:])
	md.Sum(addr[:0])

	return addr
}
