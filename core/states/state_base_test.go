
package states

import (
	"bytes"
	"testing"
)

func TestStateBase_Serialize_Deserialize(t *testing.T) {

	st := &StateBase{byte(1)}

	bf := new(bytes.Buffer)
	if err := st.Serialize(bf); err != nil {
		t.Fatalf("StateBase serialize error: %v", err)
	}

	var st2 = new(StateBase)
	if err := st2.Deserialize(bf); err != nil {
		t.Fatalf("StateBase deserialize error: %v", err)
	}
}
