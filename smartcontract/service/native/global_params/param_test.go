

package global_params

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParams_Serialize_Deserialize(t *testing.T) {
	params := Params{}
	for i := 0; i < 10; i++ {
		k := "key" + strconv.Itoa(i)
		v := "value" + strconv.Itoa(i)
		params.SetParam(Param{k, v})
	}
	bf := new(bytes.Buffer)
	if err := params.Serialize(bf); err != nil {
		t.Fatalf("params serialize error: %v", err)
	}
	deserializeParams := Params{}
	if err := deserializeParams.Deserialize(bf); err != nil {
		t.Fatalf("params deserialize error: %v", err)
	}
	for i := 0; i < 10; i++ {
		originParam := params[i]
		deseParam := deserializeParams[i]
		if originParam.Key != deseParam.Key || originParam.Value != deseParam.Value {
			t.Fatal("params deserialize error")
		}
	}
}

func TestParamNameList_Serialize_Deserialize(t *testing.T) {
	nameList := ParamNameList{}
	for i := 0; i < 3; i++ {
		nameList = append(nameList, strconv.Itoa(i))
	}
	bf := new(bytes.Buffer)
	err := nameList.Serialize(bf)
	assert.Nil(t, err)
	deserializeNameList := ParamNameList{}
	err = deserializeNameList.Deserialize(bf)
	assert.Nil(t, err)
	assert.Equal(t, nameList, deserializeNameList)
}
