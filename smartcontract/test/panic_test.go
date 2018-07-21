

package test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/core/types"
	. "github.com/mixbee/mixbee/smartcontract"
	neovm2 "github.com/mixbee/mixbee/smartcontract/service/neovm"
	"github.com/mixbee/mixbee/vm/neovm"
	"github.com/stretchr/testify/assert"
)

func TestRandomCodeCrash(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	var code []byte
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("code %x \n", code)
		}
	}()

	for i := 1; i < 10; i++ {
		fmt.Printf("test round:%d \n", i)
		code := make([]byte, i)
		for j := 0; j < 10; j++ {
			rand.Read(code)

			//cache := storage.NewCloneCache(testBatch)
			sc := SmartContract{
				Config:     config,
				Gas:        10000,
				CloneCache: nil,
			}
			engine, _ := sc.NewExecuteEngine(code)
			engine.Invoke()
		}
	}
}

func TestOpCodeDUP(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	var code = []byte{byte(neovm.DUP)}

	sc := SmartContract{
		Config:     config,
		Gas:        10000,
		CloneCache: nil,
	}
	engine, _ := sc.NewExecuteEngine(code)
	_, err := engine.Invoke()

	assert.NotNil(t, err)
}

func TestOpReadMemAttack(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	bf := new(bytes.Buffer)
	builder := neovm.NewParamsBuilder(bf)
	builder.Emit(neovm.SYSCALL)
	bs := bytes.NewBuffer(builder.ToArray())
	builder.EmitPushByteArray([]byte(neovm2.NATIVE_INVOKE_NAME))
	l := 0X7fffffc7 - 1
	serialization.WriteVarUint(bs, uint64(l))
	b := make([]byte, 4)
	bs.Write(b)

	sc := SmartContract{
		Config:     config,
		Gas:        100000,
		CloneCache: nil,
	}
	engine, _ := sc.NewExecuteEngine(bs.Bytes())
	_, err := engine.Invoke()

	assert.NotNil(t, err)

}
