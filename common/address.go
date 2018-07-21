

package common

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/itchyny/base58-go"
)

const ADDR_LEN = 20

type Address [ADDR_LEN]byte

var ADDRESS_EMPTY = Address{}

// ToHexString returns  hex string representation of Address
func (self *Address) ToHexString() string {
	return fmt.Sprintf("%x", ToArrayReverse(self[:]))
}

// Serialize serialize Address into io.Writer
func (self *Address) Serialize(w io.Writer) error {
	_, err := w.Write(self[:])
	return err
}

// Deserialize deserialize Address from io.Reader
func (self *Address) Deserialize(r io.Reader) error {
	_, err := io.ReadFull(r, self[:])
	if err != nil {
		return errors.New("deserialize Address error")
	}
	return nil
}

// ToBase58 returns base58 encoded address string
func (f *Address) ToBase58() string {
	data := append([]byte{23}, f[:]...)
	temp := sha256.Sum256(data)
	temps := sha256.Sum256(temp[:])
	data = append(data, temps[0:4]...)

	bi := new(big.Int).SetBytes(data).String()
	encoded, _ := base58.BitcoinEncoding.Encode([]byte(bi))
	return string(encoded)
}

// AddressParseFromBytes returns parsed Address
func AddressParseFromBytes(f []byte) (Address, error) {
	if len(f) != ADDR_LEN {
		return ADDRESS_EMPTY, errors.New("[Common]: AddressParseFromBytes err, len != 20")
	}

	var addr Address
	copy(addr[:], f)
	return addr, nil
}

// AddressParseFromHexString returns parsed Address
func AddressFromHexString(s string) (Address, error) {
	hx, err := HexToBytes(s)
	if err != nil {
		return ADDRESS_EMPTY, err
	}
	return AddressParseFromBytes(ToArrayReverse(hx))
}

// AddressFromBase58 returns Address from encoded base58 string
func AddressFromBase58(encoded string) (Address, error) {
	if encoded == "" {
		return ADDRESS_EMPTY, fmt.Errorf("invalid address")
	}
	decoded, err := base58.BitcoinEncoding.Decode([]byte(encoded))
	if err != nil {
		return ADDRESS_EMPTY, err
	}

	x, ok := new(big.Int).SetString(string(decoded), 10)
	if !ok {
		return ADDRESS_EMPTY, fmt.Errorf("invalid address")
	}

	buf := x.Bytes()
	if len(buf) != 1+ADDR_LEN+4 || buf[0] != byte(23) {
		return ADDRESS_EMPTY, errors.New("wrong encoded address")
	}
	ph, err := AddressParseFromBytes(buf[1:21])
	if err != nil {
		return ADDRESS_EMPTY, err
	}

	addr := ph.ToBase58()

	if addr != encoded {
		return ADDRESS_EMPTY, errors.New("[AddressFromBase58]: decode encoded verify failed.")
	}

	return ph, nil
}
