package crosschain

import (
	"testing"
	"fmt"
	"encoding/json"
	"encoding/hex"
)

func TestGetSubNetNode(t *testing.T) {

	//nodes := NewSubChainNetNodes()
	//nodes.RegisterNodes(1,"11")
	//nodes.RegisterNodes(1,"12")
	//nodes.RegisterNodes(1,"13")
	//
	//nodes.RegisterNodes(2,"21")
	//nodes.RegisterNodes(2,"22")
	//nodes.RegisterNodes(2,"23")
	//
	//n := nodes.GetSubNetNode(2)
	//fmt.Println("node info = ",n)
	//
	//n = nodes.GetSubNetNode(3)
	//fmt.Println("node info = ",n)

	pair := &CTXPairEntry{}
	first := &CTXEntry{}
	second := &CTXEntry{}
	pair.First = first
	pair.Second = second

	first.From = "fromA"
	first.To = "toA"
	first.FromValue = 1000
	first.ToValue = 100
	first.SeqId = "123456"
	first.Type = 12
	first.ANetWorkId = 1
	first.BNetWorkId = 2
	first.Sig = []byte("sig info")
	first.TxHash = "txhash"

	second.From = "fromB"
	second.To = "toB"
	second.FromValue = 100
	second.ToValue = 1000
	second.SeqId = "123456"
	second.Type = 12
	second.ANetWorkId = 1
	second.BNetWorkId = 2
	second.Sig = []byte("sig info")
	second.TxHash = "txhash"

	fmt.Printf("first=%#v\n", first)
	fmt.Printf("second=%#v\n", second)

	by, err := json.Marshal(pair)
	if err != nil {
		fmt.Printf("json marshal err\n", err)
		return
	}
	fmt.Printf("pair json str %s\n",by)
	hexStr := hex.EncodeToString(by)
	fmt.Printf("json hex str %s\n", hexStr)

	buf, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Printf("hex decode err\n", err)
		return
	}
	pairInfo := &CTXPairEntry{}
	err = json.Unmarshal(buf, pairInfo)
	if err != nil {
		fmt.Printf("json unmarshal err\n", err)
		return
	}
	fmt.Printf("first=%+v\n", pairInfo.First)
	fmt.Printf("second=%+v\n", pairInfo.Second)
}
