package common


import (
	"encoding/hex"
	"math/rand"
	"os"
	"io"
	"crypto/sha256"
	. "github.com/mixbee/mixbee/errors"
	"golang.org/x/crypto/ripemd160"
	"errors"

)



func ToCodeHash(code []byte) (Uint160, error) {
	temp := sha256.Sum256(code)
	md := ripemd160.New()
	io.WriteString(md, string(temp[:]))
	f := md.Sum(nil)

	hash, err := Uint160ParseFromBytes(f)
	if err != nil {
		return Uint160{}, NewDetailErr(errors.New("[Common] , ToCodeHash err."), ErrNoCode, "")
	}
	return hash, nil
}

// GetNonce returns random nonce
func GetNonce() uint64 {
	// Fixme replace with the real random number generator
	nonce := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	return nonce
}

// ToHexString convert []byte to hex string
func ToHexString(data []byte) string {
	return hex.EncodeToString(data)
}

// HexToBytes convert hex string to []byte
func HexToBytes(value string) ([]byte, error) {
	return hex.DecodeString(value)
}

func ToArrayReverse(arr []byte) []byte {
	l := len(arr)
	x := make([]byte, 0)
	for i := l - 1; i >= 0; i-- {
		x = append(x, arr[i])
	}
	return x
}

// FileExisted checks whether filename exists in filesystem
func FileExisted(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

