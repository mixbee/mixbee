

package vconfig

import (
	"bytes"
	"fmt"
	"testing"
)

func generTestData() []byte {
	nodeId := "12020298fe9f22e9df64f6bfcc1c2a14418846cffdbbf510d261bbc3fa6d47073df9a2"
	chainPeers := make([]*PeerConfig, 0)
	peerconfig := &PeerConfig{
		Index: 12,
		ID:    nodeId,
	}
	chainPeers = append(chainPeers, peerconfig)

	tests := &ChainConfig{
		Version:              1,
		View:                 12,
		N:                    4,
		C:                    3,
		BlockMsgDelay:        1000,
		HashMsgDelay:         1000,
		PeerHandshakeTimeout: 10000,
		Peers:                chainPeers,
		PosTable:             []uint32{2, 3, 1, 3, 1, 3, 2, 3, 2, 3, 2, 1, 3},
	}
	cc := new(bytes.Buffer)
	tests.Serialize(cc)
	return cc.Bytes()
}
func TestSerialize(t *testing.T) {
	res := generTestData()
	fmt.Println("serialize:", res)
}

func TestDeserialize(t *testing.T) {
	res := generTestData()
	test := &ChainConfig{}
	err := test.Deserialize(bytes.NewReader(res), len(res))
	if err != nil {
		t.Log("test failed ")
	}
	fmt.Printf("version: %d, C:%d \n", test.Version, test.C)
}
