

package states

import (
	"bytes"
	"io"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/serialization"
)

type StorageKey struct {
	ContractAddress common.Address
	Key             []byte
}

func (this *StorageKey) Serialize(w io.Writer) (int, error) {
	if err := this.ContractAddress.Serialize(w); err != nil {
		return 0, err
	}
	if err := serialization.WriteVarBytes(w, this.Key); err != nil {
		return 0, err
	}
	return 0, nil
}

func (this *StorageKey) Deserialize(r io.Reader) error {
	if err := this.ContractAddress.Deserialize(r); err != nil {
		return err
	}
	key, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	this.Key = key
	return nil
}

func (this *StorageKey) ToArray() []byte {
	b := new(bytes.Buffer)
	this.Serialize(b)
	return b.Bytes()
}
