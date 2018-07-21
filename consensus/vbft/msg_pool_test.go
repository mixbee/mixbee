

package vbft

import "testing"

func TestAddMsg(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, 1)
	block, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed :%v", err)
		return
	}
	blockproposalmsg := &blockProposalMsg{
		Block: block,
	}
	h, _ := HashMsg(blockproposalmsg)
	err = msgpool.AddMsg(blockproposalmsg, h)
	t.Logf("TestAddMsg %v", err)
}

func TestHasMsg(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, 1)
	block, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed :%v", err)
		return
	}
	blockproposalmsg := &blockProposalMsg{
		Block: block,
	}
	h, _ := HashMsg(blockproposalmsg)
	status := msgpool.HasMsg(blockproposalmsg, h)
	t.Logf("TestHasMsg: %v", status)
}

func TestGetProposalMsgs(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, 1)
	consensusmsgs := msgpool.GetProposalMsgs(1)
	t.Logf("TestGetProposalMsgs: %v", len(consensusmsgs))
}

func TestGetEndorsementsMsgs(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, 1)
	consensusmsgs := msgpool.GetEndorsementsMsgs(1)
	t.Logf("TestGetEndorsementsMsgs: %v", len(consensusmsgs))
}

func TestGetCommitMsgs(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, 1)
	consensusmsgs := msgpool.GetCommitMsgs(1)
	t.Logf("TestGetCommitMsgs: %v", len(consensusmsgs))
}

func TestOnBlockSealed(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
		return
	}
	blockproposalmsg := &blockProposalMsg{
		Block: blk,
	}
	h, _ := HashMsg(blockproposalmsg)
	server := constructServer()
	msgpool := newMsgPool(server, 1)
	t.Logf("TestOnBlockSealed,len:%v", len(msgpool.rounds))
	if !msgpool.HasMsg(blockproposalmsg, h) {
		msgpool.AddMsg(blockproposalmsg, h)
		msgpool.onBlockSealed(blockproposalmsg.GetBlockNum())
		t.Logf("TestOnBlockSealed,len:%v", len(msgpool.rounds))
	}
}
