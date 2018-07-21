

package utils

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/common/log"
	msgCommon "github.com/mixbee/mixbee/p2pserver/common"
	"github.com/mixbee/mixbee/p2pserver/net/netserver"
	"github.com/mixbee/mixbee/p2pserver/net/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testHandler(data *msgCommon.MsgPayload, p2p p2p.P2P, pid *actor.PID, args ...interface{}) error {
	log.Info("Test handler")
	return nil
}

// TestMsgRouter tests a basic function of a message router
func TestMsgRouter(t *testing.T) {
	_, pub, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	network := netserver.NewNetServer(pub)
	msgRouter := NewMsgRouter(network)
	assert.NotNil(t, msgRouter)

	msgRouter.RegisterMsgHandler("test", testHandler)
	msgRouter.UnRegisterMsgHandler("test")
	msgRouter.Start()
	msgRouter.Stop()
}
