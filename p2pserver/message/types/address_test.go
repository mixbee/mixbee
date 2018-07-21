

package types

import (
	"bytes"
	"net"
	"testing"

	comm "github.com/mixbee/mixbee/p2pserver/common"
	"github.com/stretchr/testify/assert"
)

func MessageTest(t *testing.T, msg Message) {
	p := new(bytes.Buffer)
	err := WriteMessage(p, msg)
	assert.Nil(t, err)

	demsg, err := ReadMessage(p)
	assert.Nil(t, err)

	assert.Equal(t, msg, demsg)
}

func TestAddressSerializationDeserialization(t *testing.T) {
	var msg Addr
	var addr [16]byte
	ip := net.ParseIP("192.168.0.1")
	ip.To16()
	copy(addr[:], ip[:16])
	nodeAddr := comm.PeerAddr{
		Time:          12345678,
		Services:      100,
		IpAddr:        addr,
		Port:          8080,
		ConsensusPort: 8081,
		ID:            987654321,
	}
	msg.NodeAddrs = append(msg.NodeAddrs, nodeAddr)

	MessageTest(t, &msg)
}
