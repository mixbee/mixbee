

package abi

import "strings"

const (
	NATIVE_PARAM_TYPE_BOOL      = "bool"
	NATIVE_PARAM_TYPE_BYTE      = "byte"
	NATIVE_PARAM_TYPE_INTEGER   = "int"
	NATIVE_PARAM_TYPE_STRING    = "string"
	NATIVE_PARAM_TYPE_BYTEARRAY = "bytearray"
	NATIVE_PARAM_TYPE_ARRAY     = "array"
	NATIVE_PARAM_TYPE_ADDRESS   = "address"
	NATIVE_PARAM_TYPE_UINT256   = "uint256"
	NATIVE_PARAM_TYPE_STRUCT    = "struct"
)

type NativeContractAbi struct {
	Address   string                       `json:"hash"`
	Functions []*NativeContractFunctionAbi `json:"functions"`
	Events    []*NativeContractEventAbi    `json:"events"`
}

type NativeContractFunctionAbi struct {
	Name       string                    `json:"name"`
	Parameters []*NativeContractParamAbi `json:"parameters"`
	ReturnType string                    `json:"returnType"`
}

type NativeContractParamAbi struct {
	Name    string                    `json:"name"`
	Type    string                    `json:"type"`
	SubType []*NativeContractParamAbi `json:"subType"`
}

type NativeContractEventAbi struct {
	Name       string                    `json:"name"`
	Parameters []*NativeContractParamAbi `json:"parameters"`
}

func (this *NativeContractAbi) GetFunc(name string) *NativeContractFunctionAbi {
	name = strings.ToLower(name)
	for _, funcAbi := range this.Functions {
		if strings.ToLower(funcAbi.Name) == name {
			return funcAbi
		}
	}
	return nil
}

func (this *NativeContractAbi) GetEvent(name string) *NativeContractEventAbi {
	name = strings.ToLower(name)
	for _, evtAbi := range this.Events {
		if strings.ToLower(evtAbi.Name) == name {
			return evtAbi
		}
	}
	return nil
}
