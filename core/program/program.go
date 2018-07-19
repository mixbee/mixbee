package program

import (
	"bytes"
	"github.com/mixbee/mixbee-crypto/keypair"
)

type ProgramBuilder struct {
	buffer bytes.Buffer
}


func (self *ProgramBuilder) PushPubKey(pubkey keypair.PublicKey) *ProgramBuilder {
	buf := keypair.SerializePublicKey(pubkey)
	return self.PushBytes(buf)
}

func (self *ProgramBuilder) Finish() []byte {
	return self.buffer.Bytes()
}

func (self *ProgramBuilder) PushBytes(data []byte) *ProgramBuilder  {
	if len(data) == 0 {
		panic("push data error: data is nil")
	}

	self.buffer.Write(data)
	return self
}