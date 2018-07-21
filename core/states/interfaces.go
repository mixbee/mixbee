

package states

import (
	"io"
)

type StateValue interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}
