package crosschain

import (
	"testing"
	"fmt"
)

func TestGetSubNetNode(t *testing.T) {

	nodes := NewSubChainNetNodes()
	nodes.RegisterNodes(1,"11")
	nodes.RegisterNodes(1,"12")
	nodes.RegisterNodes(1,"13")

	nodes.RegisterNodes(2,"21")
	nodes.RegisterNodes(2,"22")
	nodes.RegisterNodes(2,"23")

	n := nodes.GetSubNetNode(2)
	fmt.Println("node info = ",n)

	n = nodes.GetSubNetNode(3)
	fmt.Println("node info = ",n)
}
