package account

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"crypto/rand"
	"math/big"
	"github.com/itchyny/base58-go"
	"crypto/sha256"
	"bytes"
	"encoding/hex"
	"github.com/mixbee/mixbee/core/types"

)

const (
	VER    = 0x41
)

type Controller struct {
	ID     string `json:"id"`
	Public string `json:"publicKey,omitemtpy"`
	keypair.ProtectedKey
}

type Identity struct {
	ID      string       `json:"mixbeeid"`
	Label   string       `json:"label,omitempty"`
	Lock    bool         `json:"lock"`
	Control []Controller `json:"controls,omitempty"`
	Extra   interface{}  `json:"extra,omitempty"`
}

func CreateID(nonce []byte) (string, error) {
	hasher := ripemd160.New()
	_, err := hasher.Write(nonce)
	if err != nil {
		return "", fmt.Errorf("create ID error, %s", err)
	}
	data := hasher.Sum([]byte{VER})
	data = append(data, checksum(data)...)

	bi := new(big.Int).SetBytes(data).String()
	idstring, err := base58.BitcoinEncoding.Encode([]byte(bi))
	if err != nil {
		return "", fmt.Errorf("create ID error, %s", err)
	}

	return string(idstring), nil
}

func GenerateID() (string, error) {
	var buf [32]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return "", fmt.Errorf("generate ID error, %s", err)
	}
	return CreateID(buf[:])
}

func checksum(data []byte) []byte {
	sum := sha256.Sum256(data)
	sum = sha256.Sum256(sum[:])
	return sum[:4]
}

func VerifyID(id string) bool  {
	if len(id) < 9 {
		return false
	}
	buf, err := base58.BitcoinEncoding.Decode([]byte(id))
	if err != nil {
		return false
	}
	bi, ok := new(big.Int).SetString(string(buf), 10)
	if !ok || bi == nil {
		return false
	}
	buf = bi.Bytes()
	// 1 byte version + 20 byte hash + 4 byte checksum
	if len(buf) != 25 {
		return false
	}
	pos := len(buf) - 4
	data := buf[:pos]
	check := buf[pos:]
	sum := checksum(data)
	if !bytes.Equal(sum, check) {
		return false
	}
	return true
}

func NewIdentity(label string, keyType keypair.KeyType, param interface{}, password []byte) (*Identity, error) {
	var res Identity
	id,err := GenerateID()
	if err!=nil {
		return nil,err
	}
	pri, pub, err := keypair.GenerateKeyPair(keyType, param)
	if err != nil {
		return nil,err
	}
	// buf :=  keypair.SerializePublicKey(pub)
	// codeHash,_ := common.ToCodeHash(buf)
	// address,_ := codeHash.ToAddress()
	addr := types.AddressFromPubKey(pub)
	b58addr := addr.ToBase58()
	prot, err := keypair.EncryptPrivateKey(pri, b58addr, password)
	if err != nil {
		return nil, err
	}

	res.ID = id
	res.Label = label
	res.Lock = false
	res.Control = make([]Controller, 1)
	res.Control[0].ID = "1"
	res.Control[0].ProtectedKey = *prot
	res.Control[0].Public = hex.EncodeToString(keypair.SerializePublicKey(pub))
	return &res,nil
}


