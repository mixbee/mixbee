

package utils

import "github.com/mixbee/mixbee/common/constants"

var (
	TIME_INTERVAL     = constants.UNBOUND_TIME_INTERVAL
	GENERATION_AMOUNT = constants.UNBOUND_GENERATION_AMOUNT
)

// startOffset : start timestamp offset from genesis block
// endOffset :  end timestamp offset from genesis block
func CalcUnbindOng(balance uint64, startOffset, endOffset uint32) uint64 {
	var amount uint64 = 0
	if startOffset >= endOffset {
		return 0
	}
	if startOffset < constants.UNBOUND_DEADLINE {
		ustart := startOffset / TIME_INTERVAL
		istart := startOffset % TIME_INTERVAL
		if endOffset >= constants.UNBOUND_DEADLINE {
			endOffset = constants.UNBOUND_DEADLINE
		}
		uend := endOffset / TIME_INTERVAL
		iend := endOffset % TIME_INTERVAL
		for ustart < uend {
			amount += uint64(TIME_INTERVAL-istart) * GENERATION_AMOUNT[ustart]
			ustart++
			istart = 0
		}
		amount += uint64(iend-istart) * GENERATION_AMOUNT[ustart]
	}

	return uint64(amount) * balance
}
