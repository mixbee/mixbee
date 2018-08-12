package crosschain

import (
	"testing"
	"encoding/json"
	"fmt"
	"github.com/mixbee/mixbee/common/log"
)

func TestGetAllNodes(t *testing.T) {

	log.InitLog(5, log.PATH, log.Stdout)

	nodes := NewVerifyNodes()
	nodes.RegisterNodes("11", "host1")
	nodes.RegisterNodes("22","host2")

	infos := nodes.GetNodes()

	result,err := json.Marshal(infos)
	if err != nil {
		fmt.Printf("cross chain actor VerifyNodes json marshal err %s\n",err)
	}
	fmt.Printf("result %s\n",string(result))

	ll := []byte("[{\"publicKey\":\"034bdb3631ada6d57659f544acb34e429e09d721881894e42192df6912d8f63e83\",\"host\":\"http://192.168.3.4:20336\",\"timestamp\":1533890907}]")
	var list []*CrossChainVerifyNode
	err = json.Unmarshal(ll,&list)
	if err != nil {
		fmt.Printf("cross chain actor VerifyNodes json inmarshal err %s\n",err)
	}
	fmt.Printf("result %+v\n",list[0])
	buf := string(33)
	fmt.Printf("%x\n",buf)
}
