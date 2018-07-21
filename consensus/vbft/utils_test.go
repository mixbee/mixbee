

package vbft

import (
	"testing"

	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/common"
)

func HashBlock(blk *Block) (common.Uint256, error) {
	return blk.Block.Hash(), nil
}

func TestSignMsg(t *testing.T) {
	acc := account.NewAccount("SHA256withECDSA")
	if acc == nil {
		t.Error("GetDefaultAccount error: acc is nil")
		return
	}
	msg := constructProposalMsgTest(acc)
	_, err := SignMsg(acc, msg)
	if err != nil {
		t.Errorf("TestSignMsg Failed: %v", err)
		return
	}
	t.Log("TestSignMsg succ")
}

func TestHashBlock(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	hash, _ := HashBlock(blk)
	t.Logf("TestHashBlock: %v", hash)
}

func TestHashMsg(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
		return
	}
	blockproposalmsg := &blockProposalMsg{
		Block: blk,
	}
	uint256, err := HashMsg(blockproposalmsg)
	if err != nil {
		t.Errorf("TestHashMsg failed: %v", err)
		return
	}
	t.Logf("TestHashMsg succ: %v\n", uint256)
}

func TestVrf(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	vrfvalue := getParticipantSelectionSeed(blk)
	if len(vrfvalue) == 0 {
		t.Errorf("TestVrf failed:")
		return
	}
	t.Log("TestVrf succ")
}
