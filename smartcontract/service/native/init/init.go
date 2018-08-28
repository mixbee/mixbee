

package init

import (
	"bytes"
	"math"
	"math/big"

	"github.com/mixbee/mixbee/common"
	invoke "github.com/mixbee/mixbee/core/utils"
	"github.com/mixbee/mixbee/smartcontract/service/native/auth"
	params "github.com/mixbee/mixbee/smartcontract/service/native/global_params"
	"github.com/mixbee/mixbee/smartcontract/service/native/governance"
	"github.com/mixbee/mixbee/smartcontract/service/native/mbg"
	"github.com/mixbee/mixbee/smartcontract/service/native/crosspairevidence"
	"github.com/mixbee/mixbee/smartcontract/service/native/mbc"
	"github.com/mixbee/mixbee/smartcontract/service/native/mixid"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/smartcontract/service/neovm"
	vm "github.com/mixbee/mixbee/vm/neovm"
	"github.com/mixbee/mixbee/smartcontract/service/native/mixtest"
	"github.com/mixbee/mixbee/smartcontract/service/native/crosschaintx"
	"github.com/mixbee/mixbee/smartcontract/service/native/crossverifynode"
)

var (
	COMMIT_DPOS_BYTES = InitBytes(utils.GovernanceContractAddress, governance.COMMIT_DPOS)
)

func init() {
	mbg.InitMbg()
	mbc.InitMbc()
	params.InitGlobalParams()
	mixid.Init()
	auth.Init()
	governance.InitGovernance()

	mixtest.InitMixTest()
	crosschaintx.InitCrossChainTx()
	crosspairevidence.InitCrossPairEvidence()
	crossverifynode.InitCrossChainVerifyNode()
}

func InitBytes(addr common.Address, method string) []byte {

	bf := new(bytes.Buffer)
	builder := vm.NewParamsBuilder(bf)
	builder.EmitPushByteArray([]byte{})
	builder.EmitPushByteArray([]byte(method))
	builder.EmitPushByteArray(addr[:])
	builder.EmitPushInteger(big.NewInt(0))
	builder.Emit(vm.SYSCALL)
	builder.EmitPushByteArray([]byte(neovm.NATIVE_INVOKE_NAME))

	tx := invoke.NewInvokeTransaction(builder.ToArray())
	tx.GasLimit = math.MaxUint64
	return bf.Bytes()
}
