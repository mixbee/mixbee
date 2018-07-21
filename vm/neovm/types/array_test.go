

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelfArray(t *testing.T) {
	a := NewArray(nil)
	b := NewArray([]StackItems{a})
	a.Add(b)

	equ := a.Equals(b)
	assert.False(t, equ)
}
