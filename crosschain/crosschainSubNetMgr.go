package crosschain

import (
	"sync"
	"time"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/config"
	"math/rand"
)

type SubChainNetNodes struct {
	SubNetNodes map[uint32]map[string]uint32
	sync.RWMutex
}

func NewSubChainNetNodes() *SubChainNetNodes {
	nodes := &SubChainNetNodes{}
	nodes.SubNetNodes = make(map[uint32]map[string]uint32)
	return nodes
}

func (this *SubChainNetNodes) Init() {
	go func() {
		ticker := time.NewTicker(config.DEFAULT_CROSS_CHAIN_VERIFY_PING_TIME * time.Second)
		for {
			select {
			case <-ticker.C:
				deleteTimeOutSubNetNode(this)
			}
		}
	}()
}

func deleteTimeOutSubNetNode(nodes *SubChainNetNodes) {
	if len(nodes.SubNetNodes) == 0 {
		return
	}
	nodes.Lock()
	defer nodes.Unlock()
	for k, v := range nodes.SubNetNodes {

		for nk, nv := range v {
			if uint32(time.Now().Unix())-nv > config.DEFAULT_CROSS_CHAIN_VERIFY_PING_TIMEOUT {
				log.Infof("cross chain sub net node timeout %s delete", nk)
				delete(v, nk)
			}
		}

		if len(v) == 0 {
			log.Infof("cross chain sub net node timeout %s delete", k)
			delete(nodes.SubNetNodes, k)
		}
	}
}

func (this *SubChainNetNodes) RegisterNodes(netId uint32, host string) {

	this.Lock()
	defer this.Unlock()
	log.Infof("cross chain sub net node register netId=%d,host=%s,time=%d", netId, host, time.Now().Unix())

	info, ok := this.SubNetNodes[netId]
	if !ok {
		info = make(map[string]uint32)
		this.SubNetNodes[netId] = info
	}

	info[host] = uint32(time.Now().Unix())
}

func (this *SubChainNetNodes) GetSubNetNode(netId uint32) string {

	this.RLock()
	defer this.RUnlock()
	info, ok := this.SubNetNodes[netId]
	if !ok {
		return ""
	}

	if len(info) == 0 {
		return ""
	}

	var hosts []string
	for k := range info {
		hosts = append(hosts, k)
	}

	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(hosts))
	return hosts[index]
}

func (this *SubChainNetNodes) CheckNetId(ids []uint32) bool {
	this.RLock()
	defer this.RUnlock()
	if len(ids) == 0 {
		return false
	}
	for _, value := range ids {
		info, ok := this.SubNetNodes[value]
		if !ok {
			return false
		}
		if len(info) == 0 {
			return false
		}
	}
	return true
}
