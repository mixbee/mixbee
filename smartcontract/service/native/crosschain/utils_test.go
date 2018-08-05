package crosschain

import (
	"testing"
	"github.com/mixbee/mixbee/common"
	"github.com/stretchr/testify/assert"
)

func TestGetSeqId(t *testing.T) {


	from,_ := common.AddressFromBase58("AegknyPTP452dyMyxMBgcc9mZJthcY6Muz")
	to,_ := common.AddressFromBase58("AV8JdaKEZezgrpq5D4j3pKqJaMw5xCzd4n")
	state1 := &CrossChainState{
		From:from,
		To:to,
		AValue:1000,
		BValue:100,
		AChainId:3,
		BChainId:5,
		Type:1,
		Timestamp:100000001,
	}

	state2 := &CrossChainState{
		From:to,
		To:from,
		AValue:100,
		BValue:1000,
		AChainId:5,
		BChainId:3,
		Type:1,
		Timestamp:100000000,
	}

	s1 := GetSeqId(state1)
	s2 := GetSeqId(state2)

	assert.Equal(t,s1,s2)

}
