
package types

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
	"encoding/hex"
)

func TestAddressFromBookkeepers(t *testing.T) {
	_, pubKey1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey2, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey3, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	pubkeys := []keypair.PublicKey{pubKey1, pubKey2, pubKey3}

	addr, _ := AddressFromBookkeepers(pubkeys)
	addr2, _ := AddressFromMultiPubKeys(pubkeys, 3)
	assert.Equal(t, addr, addr2)

	pubkeys = []keypair.PublicKey{pubKey3, pubKey2, pubKey1}
	addr3, _ := AddressFromMultiPubKeys(pubkeys, 3)

	assert.Equal(t, addr2, addr3)
}

func TestAddressFromPubKey(t *testing.T)  {
	//  AddressFromPubKey(pubkey keypair.PublicKey)
	_, pubKey1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	add := AddressFromPubKey(pubKey1)
	// address := types.AddressFromPubKey(pub)
	t.Log("add", add)
}

func TestAddressFromMultiPubKeys(t *testing.T)  {
	// AddressFromMultiPubKeys(pubkeys []keypair.PublicKey, m int)
	//_, pubKey1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	//_, pubKey2, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	//_, pubKey3, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	//_, pubKey4, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)

	//pubkeys := []keypair.PublicKey{pubKey1, pubKey2, pubKey3, pubKey4}
	//
	//add1 := AddressFromPubKey(pubKey1)
	//t.Log("add1", add1.ToBase58())
	//
	//add2 := AddressFromPubKey(pubKey2)
	//t.Log("add2", add2.ToBase58())
	//
	//add3 := AddressFromPubKey(pubKey3)
	//t.Log("add3", add3.ToBase58())
	//
	//m := (5*len(pubkeys) + 6) / 7
	//t.Log("m", m)
	//
	//add,_ := AddressFromMultiPubKeys(pubkeys, m)
	//t.Log("add",add.ToBase58())

	a,_ := hex.DecodeString("5ab93c901928d90f790383f76c8ac8aa541194e4")
	fmt.Printf("%x",AddressFromVmCode(a))
}