
package mixid

import (
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
)

func Init() {
	native.Contracts[utils.MixIDContractAddress] = RegisterIDContract
}

func RegisterIDContract(srvc *native.NativeService) {
	srvc.Register("regIDWithPublicKey", regIdWithPublicKey)
	srvc.Register("addKey", addKey)
	srvc.Register("removeKey", removeKey)
	srvc.Register("addRecovery", addRecovery)
	srvc.Register("changeRecovery", changeRecovery)
	srvc.Register("regIDWithAttributes", regIdWithAttributes)
	srvc.Register("addAttributes", addAttributes)
	srvc.Register("removeAttribute", removeAttribute)
	srvc.Register("verifySignature", verifySignature)
	srvc.Register("getPublicKeys", GetPublicKeys)
	srvc.Register("getKeyState", GetKeyState)
	srvc.Register("getAttributes", GetAttributes)
	srvc.Register("getDDO", GetDDO)
	return
}
