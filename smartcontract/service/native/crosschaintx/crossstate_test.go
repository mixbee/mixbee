package crosschaintx

import (
	"testing"
	"bytes"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func Test (t *testing.T) {


	seqs := NewCrossSeqIds()
	seqs.SeqIds = append(seqs.SeqIds,"11")
	seqs.SeqIds = append(seqs.SeqIds,"22")
	fmt.Printf("%v\n",seqs.SeqIds)
	var bb []byte
	buf := bytes.NewBuffer(bb)
	err := seqs.Serialize(buf)
	assert.Empty(t,err)


	seq22 := NewCrossSeqIds()
	err = seq22.Deserialize(buf)
	assert.Empty(t,err)
	fmt.Printf("%v\n",seq22.SeqIds)

}
