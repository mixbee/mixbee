package states

import (
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/asset"
	"github.com/mixbee/mixbee-crypto/keypair"
	"io"
	"github.com/mixbee/mixbee/common/serialization"
	"bytes"
)
type AssetState struct {
	StateBase
	AssetId common.Uint256
	AssetType asset.AssetType
	Name string
	Amount common.Fixed64
	Available common.Fixed64
	Precision byte
	FeeMode byte
	Fee common.Fixed64
	FeeAddress *common.Uint256
	Owner keypair.PublicKey
	Admin common.Uint256
	Issuer common.Uint256
	Expiration uint32
	IsFrozen bool
}

func (assetState *AssetState)Serialize(w io.Writer) error  {
	assetState.StateBase.Serialize(w)
	assetState.AssetId.Serialize(w)
	serialization.WriteUint8(w, byte(assetState.AssetType))
	assetState.Amount.Serialize(w)
	assetState.Available.Serialize(w)
	serialization.WriteVarBytes(w, []byte{assetState.Precision})
	serialization.WriteVarBytes(w, []byte{assetState.FeeMode})
	assetState.Fee.Serialize(w)
	assetState.FeeAddress.Serialize(w)

	buf := keypair.SerializePublicKey(assetState.Owner)
	serialization.WriteVarBytes(w, buf)

	assetState.Admin.Serialize(w)
	assetState.Issuer.Serialize(w)
	serialization.WriteUint32(w, assetState.Expiration)
	serialization.WriteBool(w, assetState.IsFrozen)
	return nil
}

func (assetState *AssetState)Deserialize(r io.Reader) error {
	u256 := new(common.Uint256)
	f := new(common.Fixed64)

	stateBase := new(StateBase)
	err := stateBase.Deserialize(r)
	if err != nil {
		return err
	}
	assetState.StateBase = *stateBase

	err = u256.Deserialize(r)
	if err != nil{
		return nil
	}
	assetState.AssetId = *u256

	val,err := serialization.ReadBytes(r, 1)
	if err != nil {
		return err
	}
	assetState.AssetType = asset.AssetType(val[0])


	name, err := serialization.ReadString(r)
	if err != nil{
		return nil

	}
	assetState.Name = name

	err = f.Deserialize(r)
	if err != nil{
		return nil

	}
	assetState.Amount = *f

	err = f.Deserialize(r)
	if err != nil{
		return nil
	}
	assetState.Available = *f

	precisions, err := serialization.ReadVarBytes(r)
	if err != nil {
		return nil
	}
	assetState.Precision = precisions[0]

	feeModes, err := serialization.ReadVarBytes(r)
	if err != nil {
		return nil
	}
	assetState.FeeMode = feeModes[0]

	err = f.Deserialize(r)
	if err != nil{
		return nil
	}
	assetState.Fee = *f
	err = u256.Deserialize(r)
	if err != nil{
		return nil
	}
	assetState.FeeAddress = u256


	buf, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	pk, err := keypair.DeserializePublicKey(buf)
	if err != nil {
		return err
	}
	assetState.Owner = pk



	err = u256.Deserialize(r)
	if err != nil {
		return nil
	}
	assetState.Admin = *u256
	err = u256.Deserialize(r)
	if err != nil {
		return nil
	}
	assetState.Issuer = *u256
	expiration, err := serialization.ReadUint32(r)
	if err != nil {
		return nil
	}
	assetState.Expiration = expiration
	isFrozon, err := serialization.ReadBool(r)
	if err != nil {
		return nil
	}
	assetState.IsFrozen = isFrozon
	return nil
}

func (assetState *AssetState)ToArray() []byte  {
	b := new(bytes.Buffer)
	assetState.Serialize(b)
	return b.Bytes()
}