

package states

import (
	"io"

	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
)

type StateBase struct {
	StateVersion byte
}

func (this *StateBase) Serialize(w io.Writer) error {
	serialization.WriteByte(w, this.StateVersion)
	return nil
}

func (this *StateBase) Deserialize(r io.Reader) error {
	stateVersion, err := serialization.ReadByte(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StateBase], StateBase Deserialize failed.")
	}
	this.StateVersion = stateVersion
	return nil
}
