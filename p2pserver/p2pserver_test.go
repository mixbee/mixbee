

package p2pserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/p2pserver/common"
)

var key keypair.PublicKey
var acct *account.Account

func init() {
	log.Init(log.Stdout)
	fmt.Println("Start test the netserver...")
	acct = account.NewAccount("SHA256withECDSA")
	key = acct.PubKey()

}
func TestNewP2PServer(t *testing.T) {
	log.Init(log.Stdout)
	fmt.Println("Start test new p2pserver...")

	p2p, err := NewServer(acct)
	if err != nil {
		t.Fatalf("TestP2PActorServer: p2pserver NewServer error %s", err)
	}
	//false because the ledger actor not running
	p2p.Start(false)
	defer p2p.Stop()

	if p2p.GetVersion() != common.PROTOCOL_VERSION {
		t.Error("TestNewP2PServer p2p version error", p2p.GetVersion())
	}

	var id uint64
	k := keypair.SerializePublicKey(key)
	err = binary.Read(bytes.NewBuffer(k[:8]), binary.LittleEndian, &(id))
	if err != nil {
		t.Error(err)
	}

	if p2p.GetID() != id {
		t.Error("TestNewP2PServer p2p id error")
	}
	if p2p.GetVersion() != common.PROTOCOL_VERSION {
		t.Error("TestNewP2PServer p2p version error")
	}
	sync, cons := p2p.GetPort()
	if sync != 20338 {
		t.Error("TestNewP2PServer sync port error")
	}

	if cons != 20339 {
		t.Error("TestNewP2PServer consensus port error")
	}
	go p2p.WaitForSyncBlkFinish()
	<-time.After(time.Second * common.KEEPALIVE_TIMEOUT)
}
