package crosschain

import (
	"sync"
	"time"
	"encoding/hex"
	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"strconv"
	p2ptypes "github.com/mixbee/mixbee/p2pserver/message/types"
	"github.com/mixbee/mixbee/smartcontract/service/native/crossverifynode"
)

type CrossChainVerifyNode struct {
	PublicKey string `json:"publicKey"`
	Host      string `json:"host"`
	TimeStamp uint32 `json:"timestamp"`
	Status    uint64 `json:"status"`
	Deposit   uint64 `json:"deposit"`
}

type VerifyNodes struct {
	VerifyerNodes map[string]*CrossChainVerifyNode
	acc           *account.Account
	sync.RWMutex
}

func NewVerifyNodes() *VerifyNodes {
	nodes := &VerifyNodes{}
	nodes.VerifyerNodes = make(map[string]*CrossChainVerifyNode)
	return nodes
}

func (this *VerifyNodes) Init(acc *account.Account, p2pPid *actor.PID) {
	go func() {
		this.acc = acc
		bb, host := getVerifyNodeMetaInfo(acc)
		this.RegisterNodes(hex.EncodeToString(bb), host)
		ticker := time.NewTicker(config.DEFAULT_CROSS_CHAIN_VERIFY_PING_TIME * time.Second)
		for {
			select {
			case <-ticker.C:
				deleteTimeOutVerifyNode(this)
				pushNodeToOtherVerifyNode(acc, this, p2pPid)
			}
		}
	}()
}

func deleteTimeOutVerifyNode(nodes *VerifyNodes) {
	if len(nodes.VerifyerNodes) == 0 {
		return
	}
	nodes.Lock()
	defer nodes.Unlock()
	for k, v := range nodes.VerifyerNodes {
		if uint32(time.Now().Unix())-v.TimeStamp > config.DEFAULT_CROSS_CHAIN_VERIFY_PING_TIMEOUT {
			log.Infof("cross chain verify node timeout %s delete", v.PublicKey)
			delete(nodes.VerifyerNodes, k)
		}
	}
}

func pushNodeToOtherVerifyNode(acc *account.Account, verifyNodes *VerifyNodes, p2pPid *actor.PID) {
	bb, host := getVerifyNodeMetaInfo(acc)
	verifyNodes.RegisterNodes(hex.EncodeToString(bb), host)

	if config.DefConfig.Genesis.ConsensusType != config.CONSENSUS_TYPE_SOLO {
		info := &p2ptypes.CrossChainVerifyNode{
			PublicKey: hex.EncodeToString(bb),
			Host:      host,
			Type:      1,
		}
		p2pPid.Tell(info)
	}
}

func getVerifyNodeMetaInfo(acc *account.Account) ([]byte, string) {
	bb := keypair.SerializePublicKey(acc.PublicKey)
	ip, err := common.GetLocalIp()
	if err != nil {
		log.Errorf("get local ip err:%s", err.Error())
	}
	port := config.DefConfig.Rpc.HttpJsonPort
	portStr := strconv.FormatUint(uint64(port), 10)
	host := "http://" + ip + ":" + portStr
	return bb, host
}

func (this *VerifyNodes) RegisterNodes(pbk, host string) {

	this.Lock()
	defer this.Unlock()
	log.Debugf("cross chain verify node register pbk=%s,host=%s time=%v", pbk, host, time.Now().Unix())
	info := &CrossChainVerifyNode{
		PublicKey: pbk,
		Host:      host,
		TimeStamp: uint32(time.Now().Unix()),
	}

	//check verifyNode from native smartContract crossVerifyNode
	nodeInfo, err := getVerifyNodeInfoFromNative(pbk)
	if err != nil {
		log.Errorf("crossVerifyNode||RegisterNodes||error %s", err)
		return
	}
	if nodeInfo == nil {
		//register verifyNode
		txhash, err := CrossChainVerifyNodeRegister(this.acc, info, host)
		if err != nil {
			log.Errorf("crossVerifyNode||RegisterNodes||error %s", err)
			return
		}
		log.Infof("crossVerifyNode||RegisterNodes||txHash %s", txhash)
	} else {
		info.Deposit = nodeInfo.Deposit
		info.Status = nodeInfo.CurrentStatus
	}

	this.VerifyerNodes[pbk] = info
}

func (this *VerifyNodes) DeleteNodes(pbk string) {
	this.Lock()
	defer this.Unlock()
	delete(this.VerifyerNodes, pbk)
}

func (this *VerifyNodes) GetNodes() []*CrossChainVerifyNode {
	this.Lock()
	defer this.Unlock()

	var nodes []*CrossChainVerifyNode
	if len(this.VerifyerNodes) == 0 {
		return nodes
	}

	for _, v := range this.VerifyerNodes {
		if v.Status == crossverifynode.CanVerifyStatus {
			nodes = append(nodes, v)
		}
	}
	return nodes
}

func (this *VerifyNodes) GetNode(pbk string) *CrossChainVerifyNode {
	value, ok := this.VerifyerNodes[pbk]
	if !ok {
		return nil
	}
	return value
}
