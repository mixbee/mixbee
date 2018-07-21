

package test

import (
	"fmt"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/smartcontract"
	"github.com/mixbee/mixbee/vm/neovm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	byteCode := []byte{
		byte(neovm.NEWMAP),
		byte(neovm.DUP),   // dup map
		byte(neovm.PUSH0), // key (index)
		byte(neovm.PUSH0), // key (index)
		byte(neovm.SETITEM),

		byte(neovm.DUP),   // dup map
		byte(neovm.PUSH0), // key (index)
		byte(neovm.PUSH1), // value (newItem)
		byte(neovm.SETITEM),
	}

	// pick a value out
	byteCode = append(byteCode,
		[]byte{ // extract element
			byte(neovm.DUP),   // dup map (items)
			byte(neovm.PUSH0), // key (index)

			byte(neovm.PICKITEM),
			byte(neovm.JMPIF), // dup map (items)
			0x04, 0x00,        // skip a drop?
			byte(neovm.DROP),
		}...)

	// count faults vs successful executions
	N := 1024
	faults := 0

	//dbFile := "/tmp/test"
	//os.RemoveAll(dbFile)
	//testLevelDB, err := leveldbstore.NewLevelDBStore(dbFile)
	//if err != nil {
	//	panic(err)
	//}

	for n := 0; n < N; n++ {
		// Setup Execution Environment
		//store := statestore.NewMemDatabase()
		//testBatch := statestore.NewStateStoreBatch(store, testLevelDB)
		config := &smartcontract.Config{
			Time:   10,
			Height: 10,
			Tx:     &types.Transaction{},
		}
		//cache := storage.NewCloneCache(testBatch)
		sc := smartcontract.SmartContract{
			Config:     config,
			Gas:        100,
			CloneCache: nil,
		}
		engine, err := sc.NewExecuteEngine(byteCode)

		_, err = engine.Invoke()
		if err != nil {
			fmt.Println("err:", err)
			faults += 1
		}
	}
	assert.Equal(t, faults, 0)

}
