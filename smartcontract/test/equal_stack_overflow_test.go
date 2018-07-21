

package test

import (
	"os"
	"testing"

	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/core/types"
	. "github.com/mixbee/mixbee/smartcontract"
	"github.com/mixbee/mixbee/vm/neovm"
	"github.com/stretchr/testify/assert"
)

func TestEqualStackOverflow(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("./Log")
	}()

	code := []byte{
		byte(neovm.PUSH1),    // {1}
		byte(neovm.NEWARRAY), // {[]}
		byte(neovm.DUP),      // {[],[]}
		byte(neovm.DUP),      // {[],[],[]}
		byte(neovm.PUSH0),    // {[],[],[],0}
		byte(neovm.ROT),      // {[],[],0,[]}
		byte(neovm.SETITEM),  // {[[]]}
		byte(neovm.DUP),      // {[[]],[[]]}
		byte(neovm.EQUAL),
	}

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}
	sc := SmartContract{
		Config:     config,
		Gas:        10000,
		CloneCache: nil,
	}
	engine, _ := sc.NewExecuteEngine(code)
	_, err := engine.Invoke()

	assert.Nil(t, err)
}
