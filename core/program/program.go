package program

import (
	"github.com/mixbee/mixbee/vm/avm"
	"github.com/mixbee/mixbee-crypto/keypair"
	"bytes"
	"github.com/mixbee/mixbee/common/serialization"
)

type ProgramBuilder struct {
	buffer bytes.Buffer
}

func ProgramFromPubKey(pubkey keypair.PublicKey) []byte {
	builder := ProgramBuilder{}
	return builder.PushPubKey(pubkey).PushOpCode(avm.CHECKSIG).Finish()
}

func (self *ProgramBuilder) PushPubKey(pubkey keypair.PublicKey) *ProgramBuilder {
	buf := keypair.SerializePublicKey(pubkey)
	return self.PushBytes(buf)
}

func (self *ProgramBuilder) PushOpCode(op avm.OpCode) *ProgramBuilder {
	self.buffer.WriteByte(byte(op))
	return self
}

func (self *ProgramBuilder) Finish() []byte {
	return self.buffer.Bytes()
}


func (self *ProgramBuilder) PushBytes(data []byte) *ProgramBuilder {
	if len(data) == 0 {
		panic("push data error: data is nil")
	}

	if len(data) <= int(avm.PUSHBYTES75)+1-int(avm.PUSHBYTES1) {
		self.buffer.WriteByte(byte(len(data)) + byte(avm.PUSHBYTES1) - 1)
	} else if len(data) < 0x100 {
		self.buffer.WriteByte(byte(avm.PUSHDATA1))
		serialization.WriteUint8(&self.buffer, uint8(len(data)))
	} else if len(data) < 0x10000 {
		self.buffer.WriteByte(byte(avm.PUSHDATA2))
		serialization.WriteUint16(&self.buffer, uint16(len(data)))
	} else {
		self.buffer.WriteByte(byte(avm.PUSHDATA4))
		serialization.WriteUint32(&self.buffer, uint32(len(data)))
	}
	self.buffer.Write(data)

	return self
}