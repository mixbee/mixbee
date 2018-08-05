package common

import (
	"encoding/hex"
	"math/rand"
	"os"
	"net"
	"fmt"
)

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

func GetLocalIp() (string,error){

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "",err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(),nil
			}

		}
	}
	return "",fmt.Errorf("not get local ip")
}
