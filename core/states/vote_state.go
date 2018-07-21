
package states

import (
	"io"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/serialization"
)

type VoteState struct {
	StateBase
	PublicKeys []keypair.PublicKey
	Count      common.Fixed64
}

func (this *VoteState) Serialize(w io.Writer) error {
	err := this.StateBase.Serialize(w)
	if err != nil {
		return err
	}
	err = serialization.WriteUint32(w, uint32(len(this.PublicKeys)))
	if err != nil {
		return err
	}
	for _, v := range this.PublicKeys {
		buf := keypair.SerializePublicKey(v)
		err := serialization.WriteVarBytes(w, buf)
		if err != nil {
			return err
		}
	}

	return serialization.WriteUint64(w, uint64(this.Count))
}

func (this *VoteState) Deserialize(r io.Reader) error {
	err := this.StateBase.Deserialize(r)
	if err != nil {
		return err
	}
	n, err := serialization.ReadUint32(r)
	if err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		buf, err := serialization.ReadVarBytes(r)
		if err != nil {
			return err
		}
		pk, err := keypair.DeserializePublicKey(buf)
		if err != nil {
			return err
		}
		this.PublicKeys = append(this.PublicKeys, pk)
	}
	c, err := serialization.ReadUint64(r)
	if err != nil {
		return err
	}
	this.Count = common.Fixed64(int64(c))
	return nil
}
