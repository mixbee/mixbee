
package common

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/binary"
	"fmt"
	"strconv"
)

func TestUint256_Serialize(t *testing.T) {
	var val Uint256
	val[1] = 245
	buf := bytes.NewBuffer(nil)
	err := val.Serialize(buf)
	assert.NotNil(t, err)
}

func TestUint256_Deserialize(t *testing.T) {
	var val Uint256
	val[1] = 245
	buf := bytes.NewBuffer(nil)
	val.Serialize(buf)

	var val2 Uint256
	val2.Deserialize(buf)

	assert.Equal(t, val, val2)

	buf = bytes.NewBuffer([]byte{1, 2, 3})
	err := val2.Deserialize(buf)

	assert.NotNil(t, err)
}

func TestUint256ParseFromBytes(t *testing.T) {
	buf := []byte{1, 2, 3}

	_, err := Uint256ParseFromBytes(buf)
	t.Log("err",err)

	bs1 := make([]byte,32)
	binary.LittleEndian.PutUint32(bs1, 31415926)
	fmt.Println(bs1)

	bs2 := []byte(strconv.Itoa(31415926))
	fmt.Println(bs2)

	assert.NotNil(t, err)
}
