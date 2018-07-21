

package abi

import (
	"encoding/json"
	"fmt"
	"github.com/mixbee/mixbee/common/log"
	"io/ioutil"
	"strings"
)

var DefAbiMgr = NewAbiMgr()

type AbiMgr struct {
	Path       string
	nativeAbis map[string]*NativeContractAbi
}

func NewAbiMgr() *AbiMgr {
	return &AbiMgr{
		nativeAbis: make(map[string]*NativeContractAbi),
	}
}

func (this *AbiMgr) GetNativeAbi(address string) *NativeContractAbi {
	abi, ok := this.nativeAbis[address]
	if ok {
		return abi
	}
	return nil
}

func (this *AbiMgr) Init(path string) {
	this.Path = path
	this.loadNativeAbi()
}

func (this *AbiMgr) loadNativeAbi() {
	nativeAbiFiles, err := ioutil.ReadDir(this.Path)
	if err != nil {
		log.Errorf("AbiMgr loadNativeAbi read dir:./native error:%s", err)
		return
	}
	for _, nativeAbiFile := range nativeAbiFiles {
		fileName := nativeAbiFile.Name()
		if nativeAbiFile.IsDir() {
			continue
		}
		if !strings.HasSuffix(fileName, ".json") {
			continue
		}
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", this.Path, fileName))
		if err != nil {
			log.Errorf("AbiMgr loadNativeAbi name:%s error:%s", fileName, err)
			continue
		}
		nativeAbi := &NativeContractAbi{}
		err = json.Unmarshal(data, nativeAbi)
		if err != nil {
			log.Errorf("AbiMgr loadNativeAbi name:%s error:%s", fileName, err)
			continue
		}
		this.nativeAbis[nativeAbi.Address] = nativeAbi
		log.Infof("Native contract name:%s address:%s abi load success", fileName, nativeAbi.Address)
	}
}
