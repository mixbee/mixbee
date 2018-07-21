

package common

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLimitedWriter_Write(t *testing.T) {
	bf := bytes.NewBuffer(nil)
	writer := NewLimitedWriter(bf, 5)
	_, err := writer.Write([]byte{1, 2, 3})
	assert.Nil(t, err)
	assert.Equal(t, bf.Bytes(), []byte{1, 2, 3})
	_, err = writer.Write([]byte{4, 5})
	assert.Nil(t, err)

	_, err = writer.Write([]byte{6})
	assert.Equal(t, err, ErrWriteExceedLimitedCount)
}
