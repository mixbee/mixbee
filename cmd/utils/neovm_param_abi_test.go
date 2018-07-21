
package utils

import (
	"fmt"
	"testing"
)

func TestParseNeovmFunc(t *testing.T) {
	var testNeovmAbi = `{
  "hash": "0xe827bf96529b5780ad0702757b8bad315e2bb8ce",
  "entrypoint": "Main",
  "functions": [
    {
      "name": "Main",
      "parameters": [
        {
          "name": "operation",
          "type": "String"
        },
        {
          "name": "args",
          "type": "Array"
        }
      ],
      "returntype": "Any"
    },
    {
      "name": "Add",
      "parameters": [
        {
          "name": "a",
          "type": "Integer"
        },
        {
          "name": "b",
          "type": "Integer"
        }
      ],
      "returntype": "Integer"
    }
  ],
  "events": []
}`
	contractAbi, err := NewNeovmContractAbi([]byte(testNeovmAbi))
	if err != nil {
		t.Error("TestParseNeovmFunc NewNeovmContractAbi error:%s", err)
		return
	}
	funcAbi := contractAbi.GetFunc("Add")
	if funcAbi == nil {
		t.Error("TestParseNeovmFunc cannot find func abi")
		return
	}

	params, err := ParseNeovmFunc([]string{"12", "34"}, funcAbi)
	if err != nil {
		t.Error("TestParseNeovmFunc ParseNeovmFunc error:%s", err)
		return
	}
	fmt.Printf("TestParseNeovmFunc %v\n", params)
}
