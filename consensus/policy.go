

package consensus

import (
	"github.com/mixbee/mixbee/common"
)

type Policy struct {
	PolicyLevel PolicyLevel
	List        []common.Address
}

func NewPolicy() *Policy {
	return &Policy{}
}

func (p *Policy) Refresh() {
	//TODO: Refresh
}

var DefaultPolicy *Policy

func InitPolicy() {
	DefaultPolicy := NewPolicy()
	DefaultPolicy.Refresh()
}
