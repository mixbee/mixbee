

package states

import (
	"bytes"
	"io"

	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
)

type StorageItem struct {
	StateBase
	Value []byte
}

func (this *StorageItem) Serialize(w io.Writer) error {
	this.StateBase.Serialize(w)
	serialization.WriteVarBytes(w, this.Value)
	return nil
}

func (this *StorageItem) Deserialize(r io.Reader) error {
	err := this.StateBase.Deserialize(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageItem], StateBase Deserialize failed.")
	}
	value, err := serialization.ReadVarBytes(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageItem], Value Deserialize failed.")
	}
	this.Value = value
	return nil
}

func (storageItem *StorageItem) ToArray() []byte {
	b := new(bytes.Buffer)
	storageItem.Serialize(b)
	return b.Bytes()
}
