package crosschain

import 	(
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

func TestGetTxStateByHash(t *testing.T) {
	addr := "http://localhost:20336"
	hash := "0f513ed6dfe56d632f5aebe436072ce607c4ea41cfd57668f56201f3ebfdeaf5"
	state,err := GetTxStateByHash(addr,hash)
	assert.Nil(t,err)
	fmt.Println("state：",state)
}

