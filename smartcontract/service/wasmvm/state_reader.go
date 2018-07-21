

package wasmvm

import (
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/vm/wasmvm/exec"
)

type WasmStateReader struct {
	serviceMap map[string]func(*exec.ExecutionEngine) (bool, error)
}

func NewWasmStateReader() *WasmStateReader {
	i := &WasmStateReader{
		serviceMap: make(map[string]func(*exec.ExecutionEngine) (bool, error)),
	}
	return i
}

func (i *WasmStateReader) Register(name string, handler func(*exec.ExecutionEngine) (bool, error)) bool {
	if _, ok := i.serviceMap[name]; ok {
		return false
	}
	i.serviceMap[name] = handler
	return true
}

func (i *WasmStateReader) Invoke(methodName string, engine *exec.ExecutionEngine) (bool, error) {

	if v, ok := i.serviceMap[methodName]; ok {
		return v(engine)
	}
	return true, errors.NewErr("Not supported method:" + methodName)
}

func (i *WasmStateReader) MergeMap(mMap map[string]func(*exec.ExecutionEngine) (bool, error)) bool {

	for k, v := range mMap {
		if _, ok := i.serviceMap[k]; !ok {
			i.serviceMap[k] = v
		}
	}
	return true
}

func (i *WasmStateReader) GetServiceMap() map[string]func(*exec.ExecutionEngine) (bool, error) {
	return i.serviceMap
}

func (i *WasmStateReader) Exists(name string) bool {
	_, ok := i.serviceMap[name]
	return ok
}
