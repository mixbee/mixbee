

package actor

import (
	"errors"
	"time"

	"github.com/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/common/log"
	ac "github.com/mixbee/mixbee/p2pserver/actor/server"
	"github.com/mixbee/mixbee/p2pserver/common"
)

var netServerPid *actor.PID

func SetNetServerPID(actr *actor.PID) {
	netServerPid = actr
}

func Xmit(msg interface{}) error {
	if netServerPid == nil {
		return nil
	}
	netServerPid.Tell(msg)
	return nil
}

func GetConnectionCnt() (uint32, error) {
	if netServerPid == nil {
		return 1, nil
	}
	future := netServerPid.RequestFuture(&ac.GetConnectionCntReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	r, ok := result.(*ac.GetConnectionCntRsp)
	if !ok {
		return 0, errors.New("fail")
	}
	return r.Cnt, nil
}

func GetNeighborAddrs() []common.PeerAddr {
	if netServerPid == nil {
		return []common.PeerAddr{}
	}
	future := netServerPid.RequestFuture(&ac.GetNeighborAddrsReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil
	}
	r, ok := result.(*ac.GetNeighborAddrsRsp)
	if !ok {
		return nil
	}
	return r.Addrs
}

func GetConnectionState() (uint32, error) {
	if netServerPid == nil {
		return 0, nil
	}
	future := netServerPid.RequestFuture(&ac.GetConnectionStateReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	r, ok := result.(*ac.GetConnectionStateRsp)
	if !ok {
		return 0, errors.New("fail")
	}
	return r.State, nil
}

func GetNodeTime() (int64, error) {
	if netServerPid == nil {
		return 0, nil
	}
	future := netServerPid.RequestFuture(&ac.GetTimeReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	r, ok := result.(*ac.GetTimeRsp)
	if !ok {
		return 0, errors.New("fail")
	}
	return r.Time, nil
}

func GetNodePort() (uint16, error) {
	if netServerPid == nil {
		return 0, nil
	}
	future := netServerPid.RequestFuture(&ac.GetPortReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	r, ok := result.(*ac.GetPortRsp)
	if !ok {
		return 0, errors.New("fail")
	}
	return r.SyncPort, nil
}

func GetID() (uint64, error) {
	if netServerPid == nil {
		return 0, nil
	}
	future := netServerPid.RequestFuture(&ac.GetIdReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	r, ok := result.(*ac.GetIdRsp)
	if !ok {
		return 0, errors.New("fail")
	}
	return r.Id, nil
}

func GetRelayState() (bool, error) {
	if netServerPid == nil {
		return false, nil
	}
	future := netServerPid.RequestFuture(&ac.GetRelayStateReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return false, err
	}
	r, ok := result.(*ac.GetRelayStateRsp)
	if !ok {
		return false, errors.New("fail")
	}
	return r.Relay, nil
}

func GetVersion() (uint32, error) {
	if netServerPid == nil {
		return 0, nil
	}
	future := netServerPid.RequestFuture(&ac.GetVersionReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	r, ok := result.(*ac.GetVersionRsp)
	if !ok {
		return 0, errors.New("fail")
	}
	return r.Version, nil
}

func GetNodeType() (uint64, error) {
	if netServerPid == nil {
		return 0, nil
	}
	future := netServerPid.RequestFuture(&ac.GetNodeTypeReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	r, ok := result.(*ac.GetNodeTypeRsp)
	if !ok {
		return 0, errors.New("fail")
	}
	return r.NodeType, nil
}
