
package states

import (
	"bytes"
	"github.com/mixbee/mixbee/core/types"
	"testing"
)

func TestContract_Serialize_Deserialize(t *testing.T) {
	addr := types.AddressFromVmCode([]byte{1})

	c := &Contract{
		Version: 0,
		Address: addr,
		Method:  "init",
		Args:    []byte{2},
	}
	bf := new(bytes.Buffer)
	if err := c.Serialize(bf); err != nil {
		t.Fatalf("Contract serialize error: %v", err)
	}

	v := new(Contract)
	if err := v.Deserialize(bf); err != nil {
		t.Fatalf("Contract deserialize error: %v", err)
	}
}
