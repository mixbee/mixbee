

package msgpack

import (
	"time"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	ct "github.com/mixbee/mixbee/core/types"
	msgCommon "github.com/mixbee/mixbee/p2pserver/common"
	mt "github.com/mixbee/mixbee/p2pserver/message/types"
	p2pnet "github.com/mixbee/mixbee/p2pserver/net/protocol"
	p2ptypes "github.com/mixbee/mixbee/p2pserver/message/types"
)

//Peer address package
func NewAddrs(nodeAddrs []msgCommon.PeerAddr) mt.Message {
	var addr mt.Addr
	addr.NodeAddrs = nodeAddrs

	return &addr
}

//Peer address request package
func NewAddrReq() mt.Message {
	var msg mt.AddrReq
	return &msg
}

///block package
func NewBlock(bk *ct.Block) mt.Message {
	log.Debug()
	var blk mt.Block
	blk.Blk = *bk

	return &blk
}

//blk hdr package
func NewHeaders(headers []*ct.Header) mt.Message {
	var blkHdr mt.BlkHeader
	blkHdr.BlkHdr = headers

	return &blkHdr
}

//blk hdr req package
func NewHeadersReq(curHdrHash common.Uint256) mt.Message {
	var h mt.HeadersReq
	h.Len = 1
	buf := curHdrHash
	copy(h.HashEnd[:], buf[:])

	return &h
}

////Consensus info package
func NewConsensus(cp *mt.ConsensusPayload) mt.Message {
	log.Debug()
	var cons mt.Consensus
	cons.Cons = *cp

	return &cons
}

func NewVerifyNode(node *p2ptypes.CrossChainVerifyNode) mt.Message {
	return &mt.CrossVerifyNodePayload{
			Node:*node,
	}
}

//InvPayload
func NewInvPayload(invType common.InventoryType, msg []common.Uint256) *mt.InvPayload {
	return &mt.InvPayload{
		InvType: invType,
		Blk:     msg,
	}
}

//Inv request package
func NewInv(invPayload *mt.InvPayload) mt.Message {
	var inv mt.Inv
	inv.P.Blk = invPayload.Blk
	inv.P.InvType = invPayload.InvType

	return &inv
}

//NotFound package
func NewNotFound(hash common.Uint256) mt.Message {
	log.Debug()
	var notFound mt.NotFound
	notFound.Hash = hash

	return &notFound
}

//ping msg package
func NewPingMsg(height uint64) *mt.Ping {
	log.Debug()
	var ping mt.Ping
	ping.Height = uint64(height)

	return &ping
}

//pong msg package
func NewPongMsg(height uint64) *mt.Pong {
	log.Debug()
	var pong mt.Pong
	pong.Height = uint64(height)

	return &pong
}

//Transaction package
func NewTxn(txn *ct.Transaction) mt.Message {
	log.Debug()
	var trn mt.Trn
	trn.Txn = txn

	return &trn
}

//version ack package
func NewVerAck(isConsensus bool) mt.Message {
	var verAck mt.VerACK
	verAck.IsConsensus = isConsensus

	return &verAck
}

//Version package
func NewVersion(n p2pnet.P2P, isCons bool, height uint32) mt.Message {
	var version mt.Version
	version.P = mt.VersionPayload{
		Version:      n.GetVersion(),
		Services:     n.GetServices(),
		SyncPort:     n.GetSyncPort(),
		ConsPort:     n.GetConsPort(),
		Nonce:        n.GetID(),
		IsConsensus:  isCons,
		HttpInfoPort: n.GetHttpInfoPort(),
		StartHeight:  uint64(height),
		TimeStamp:    time.Now().UnixNano(),
	}

	if n.GetRelay() {
		version.P.Relay = 1
	} else {
		version.P.Relay = 0
	}
	if config.DefConfig.P2PNode.HttpInfoPort > 0 {
		version.P.Cap[msgCommon.HTTP_INFO_FLAG] = 0x01
	} else {
		version.P.Cap[msgCommon.HTTP_INFO_FLAG] = 0x00
	}
	return &version
}

//transaction request package
func NewTxnDataReq(hash common.Uint256) mt.Message {
	var dataReq mt.DataReq
	dataReq.DataType = common.TRANSACTION
	dataReq.Hash = hash

	return &dataReq
}

//block request package
func NewBlkDataReq(hash common.Uint256) mt.Message {
	var dataReq mt.DataReq
	dataReq.DataType = common.BLOCK
	dataReq.Hash = hash

	return &dataReq
}

//consensus request package
func NewConsensusDataReq(hash common.Uint256) mt.Message {
	var dataReq mt.DataReq
	dataReq.DataType = common.CONSENSUS
	dataReq.Hash = hash

	return &dataReq
}
