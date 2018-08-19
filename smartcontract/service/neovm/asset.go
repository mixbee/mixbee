package neovm

import (
	vm "github.com/mixbee/mixbee/vm/neovm"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/core/asset"
	. "github.com/mixbee/mixbee/smartcontract/errors"
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/states"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"bytes"
	"fmt"
)

// 获得该资产的 ID
func AssetGetAssetId(service *NeoVmService, engine *vm.ExecutionEngine)  error{
	d,err := vm.PopInteropInterface(engine)
	if err != nil {
		return fmt.Errorf("%v", "Get AssetState error in function AssetGetAssetId")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(engine, assetState.AssetId.ToArray())
	return nil
}

// 获得该资产的类别
func AssetGetAssetType(service *NeoVmService, engine *vm.ExecutionEngine)  error{
	d,err := vm.PopInteropInterface(engine)
	if err != nil {
		return fmt.Errorf("%v", "Get AssetState error in function AssetGetAssetType")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(engine, int(assetState.AssetType))
	return nil
}


// 获得该资产的总量
func AssetGetAmount(service *NeoVmService, engine *vm.ExecutionEngine)  error{
	d,err := vm.PopInteropInterface(engine)
	if err != nil {
		return fmt.Errorf("%v", "Get AssetState error in function AssetGetAmount")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(engine, assetState.Amount.GetData())
	return nil
}

// 获得该资产的已经发行出去的数量
func AssetGetAvailable(service *NeoVmService, engine *vm.ExecutionEngine)  error{
	d,err := vm.PopInteropInterface(engine)
	if err == nil {
		return fmt.Errorf("%v", "Get AssetState error in function AssetGetAvailable")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(engine, assetState.Available.GetData())
	return nil
}


// 获得该资产的精度（最小分割数量），单位为小数点之后的位数
func AssetGetPrecision(service *NeoVmService, engine *vm.ExecutionEngine)  error{
	d,err := vm.PopInteropInterface(engine)
	if err!=nil {
		return  fmt.Errorf("%v", "Get AssetState error in function AssetGetPrecision")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(engine, int(assetState.Precision))
	return nil
}


// 获得该资产的所有人（公钥）
func AssetGetOwner(service *NeoVmService, engine *vm.ExecutionEngine)  error{
	d,err := vm.PopInteropInterface(engine)
	if err!=nil {
		return  fmt.Errorf("%v", "Get AssetState error in function AssetGetOwner")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(engine, assetState.Owner)
	return nil
}


// 获得该资产的管理员（合约地址），有权对资产的属性（如总量，名称等）进行修改
func AssetGetAdmin(service *NeoVmService, engine *vm.ExecutionEngine)  error{
	d,err := vm.PopInteropInterface(engine)
	if err!=nil {
		return  fmt.Errorf("%v", "Get AssetState error in function AssetGetAdmin")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(engine, assetState.Admin.ToArray())
	return nil
}


// 获得该资产的发行人（合约地址），有权进行资产的发行
func AssetGetIssuer(service *NeoVmService, engine *vm.ExecutionEngine)  error{
	d,err := vm.PopInteropInterface(engine)
	if err !=nil {
		return  fmt.Errorf("%v", "Get AssetState error in function AssetGetIssuer")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(engine, assetState.Issuer.ToArray())
	return nil
}


// new 注册一种资产
func CreateAsset(service *NeoVmService, engine *vm.ExecutionEngine) error{
	// asset_type, name, amount, precision, owner, admin, issuer
	txn, _ := vm.PopInteropInterface(engine)
	tx := txn.(*types.Transaction)
	assetId := tx.Hash()
	d,err := vm.PopInt(engine)
	if err != nil {
		return err
	}
	assertType :=  asset.AssetType(d)

	name,err := vm.PopByteArray(engine)
	if err!= nil {
		return  err
	}

	amount,err := vm.PopBigInt(engine)
	if err != nil {
		return err
	}
	if amount.Int64() == 0 {
		return ERR_ASSET_AMOUNT_INVALID
	}

	precision,err := vm.PopBigInt(engine)
	if precision.Int64() > 8 {
		return ERR_ASSET_PRECISION_INVALID
	}

	ownerByte,err := vm.PopByteArray(engine)
	//owner, err := crypto.DecodePoint(ownerByte)
	owner,err := keypair.DeserializePublicKey(ownerByte)
	if err != nil {
		return err
	}

	adminByte,err := vm.PopByteArray(engine)
	admin, err := common.Uint256ParseFromBytes(adminByte)
	if err != nil {
		return err
	}

	issueByte,err := vm.PopByteArray(engine)
	issue, err := common.Uint256ParseFromBytes(issueByte)
	if err != nil {
		return err
	}

	assetState := &states.AssetState{
		AssetId:    assetId,
		AssetType:  asset.AssetType(assertType),
		Name:       string(name),
		Amount:     common.Fixed64(amount.Int64()),
		Precision:  byte(precision.Int64()),
		Admin:      admin,
		Issuer:     issue,
		Owner:      owner,
		Expiration: service.Height + 1 + 2000000,
		IsFrozen:   false,
	}
	service.CloneCache.Add(scommon.ST_AssetState, assetId.ToArray(), assetState)
	vm.PushData(engine, assetState)
	return nil
}

// new 为资产续费
func AssetRenew(service *NeoVmService, engine *vm.ExecutionEngine) error {
	data, _ := vm.PopInteropInterface(engine)
	years,_ := vm.PopInt(engine)
	at := data.(*states.AssetState)
	height := service.Height + 1
	b := new(bytes.Buffer)
	at.AssetId.Serialize(b)
	state, err  := service.CloneCache.Store.TryGet(scommon.ST_AssetState, b.Bytes())
	if err != nil {
		return  fmt.Errorf("%v", "Get AssetState error in function AssetRenew")
	}

	assetState := state.Value.(*states.AssetState)
	if assetState.Expiration < height {
		assetState.Expiration = height
	}
	assetState.Expiration += uint32(years) * 2000000
	return nil
}