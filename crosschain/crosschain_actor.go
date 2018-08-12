package crosschain

import (
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/common/log"
	p2ptypes "github.com/mixbee/mixbee/p2pserver/message/types"
	httpactor "github.com/mixbee/mixbee/http/base/actor"
)

type CrossChainActor struct {
	props  *actor.Props
	server *CTXPoolServer
}

func NewCrossChainActor(Server *CTXPoolServer) *CrossChainActor {
	return &CrossChainActor{
		server: Server,
	}
}

//start a actor called net_server
func (this *CrossChainActor) Start() (*actor.PID, error) {
	this.props = actor.FromProducer(func() actor.Actor { return this })
	p2pPid, err := actor.SpawnNamed(this.props, "cross_chain_actor")
	return p2pPid, err
}

func (this *CrossChainActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Restarting:
		log.Info("cross actor restarting")
	case *actor.Stopping:
		log.Info("cross actor stopping")
	case *actor.Stopped:
		log.Info("cross actor stopped")
	case *actor.Started:
		log.Info("cross actor started")
	case *p2ptypes.CrossVerifyNodePayload:
		log.Info("cross actor CrossVerifyNodePayload", msg)
		this.server.VerifyNodes.RegisterNodes(msg.Node.PublicKey, msg.Node.Host)
	case *p2ptypes.CrossSubNetNodePayload:
		log.Info("cross actor CrossSubNetNodePayload", msg)
		this.server.SubNetNodesMgr.RegisterNodes(msg.NetId, msg.Host)
	case *httpactor.CrossSubNetNodeReq:
		log.Info("cross actor CrossSubNetNodeReq", msg)
		this.server.P2pPid.Tell(&p2ptypes.CrossSubNetNodePayload{
			NetId: msg.NetId,
			Host:  msg.Host,
		})
		this.server.SubNetNodesMgr.RegisterNodes(msg.NetId, msg.Host)
	case *httpactor.GetAllVerifyNodeReq:
		log.Info("cross actor GetAllVerifyNodeReq", msg)
		nodes := this.server.VerifyNodes.GetNodes()
		ctx.Sender().Request(nodes, ctx.Self())
	case *httpactor.CheckSubNetId:
		log.Info("cross actor CheckSubNetId", msg)
		ok := this.server.SubNetNodesMgr.CheckNetId(msg.NetIds)
		ctx.Sender().Request(ok, ctx.Self())
	case *httpactor.PushCrossChainTxRsq:
		log.Info("cross actor PushCrossChainTxRsq", msg)
		this.server.PushCtxToPool(msg)
	case *p2ptypes.CrossChainTxInfoPayload:
		log.Warnf("cross chain actor CrossChainTxInfoPayload %v", msg)
		entry := &CTXEntry{
			From:       msg.From,
			To:         msg.To,
			FromValue:  msg.FromValue,
			ToValue:    msg.ToValue,
			TxHash:     msg.TxHash,
			ANetWorkId: msg.ANetWorkId,
			BNetWorkId: msg.BNetWorkId,
			Type:       msg.Type,
			SeqId:      msg.SeqId,
			Sig:        msg.Sig,
			Pubk:       msg.Pubk,
			TimeStamp:  msg.TimeStamp,
			Nonce:      msg.Nonce,
		}
		this.server.pairTxPending.push(entry)
	case *p2ptypes.CrossChainTxCompletedPayload:
		this.server.CrossTxCompletedHandler(msg)
	default:
		log.Warnf("cross actor now handle msg=%#v", msg)
	}
}
