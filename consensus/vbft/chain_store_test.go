

package vbft

import (
	"fmt"
	"os"
	"testing"

	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/core/ledger"
)

func newChainStore() *ChainStore {
	log.Init(log.PATH, log.Stdout)
	var err error
	acct := account.NewAccount("SHA256withECDSA")
	if acct == nil {
		fmt.Println("GetDefaultAccount error: acc is nil")
		os.Exit(1)
	}

	ledger.DefLedger, err = ledger.NewLedger(config.DEFAULT_DATA_DIR)
	if err != nil {
		log.Fatalf("NewLedger error %s", err)
		os.Exit(1)
	}
	store, err := OpenBlockStore(ledger.DefLedger)
	if err != nil {
		fmt.Printf("openblockstore failed: %v\n", err)
		return nil
	}
	return store
}

func TestGetChainedBlockNum(t *testing.T) {
	chainstore := newChainStore()
	if chainstore == nil {
		t.Error("newChainStrore error")
		return
	}
	blocknum := chainstore.GetChainedBlockNum()
	t.Logf("TestGetChainedBlockNum :%d", blocknum)
}
