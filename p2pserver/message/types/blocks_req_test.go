

package types

import (
	"testing"

	cm "github.com/mixbee/mixbee/common"
)

func TestBlkReqSerializationDeserialization(t *testing.T) {
	var msg BlocksReq
	msg.HeaderHashCount = 1

	hashstr := "8932da73f52b1e22f30c609988ed1f693b6144f74fed9a2a20869afa7abfdf5e"
	bhash, _ := cm.HexToBytes(hashstr)
	copy(msg.HashStart[:], bhash)

	MessageTest(t, &msg)
}
