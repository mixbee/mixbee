
package payload

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployCode_Serialize(t *testing.T) {
	deploy := DeployCode{
		Code: []byte{1, 2, 3},
	}

	buf := bytes.NewBuffer(nil)
	deploy.Serialize(buf)
	bs := buf.Bytes()
	var deploy2 DeployCode
	deploy2.Deserialize(buf)
	assert.Equal(t, deploy2, deploy)

	buf = bytes.NewBuffer(bs[:len(bs)-1])
	err := deploy2.Deserialize(buf)
	assert.NotNil(t, err)
}
