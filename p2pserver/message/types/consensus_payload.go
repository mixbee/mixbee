

package types

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/core/signature"
	"github.com/mixbee/mixbee/errors"
)

type ConsensusPayload struct {
	Version         uint32
	PrevHash        common.Uint256
	Height          uint32
	BookkeeperIndex uint16
	Timestamp       uint32
	Data            []byte
	Owner           keypair.PublicKey
	Signature       []byte
	PeerId          uint64
	hash            common.Uint256
}

//get the consensus payload hash
func (this *ConsensusPayload) Hash() common.Uint256 {
	return common.Uint256{}
}

//Check whether header is correct
func (this *ConsensusPayload) Verify() error {
	buf := new(bytes.Buffer)
	err := this.SerializeUnsigned(buf)
	if err != nil {
		return err
	}
	err = signature.Verify(this.Owner, buf.Bytes(), this.Signature)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetVerifyFail, fmt.Sprintf("signature verify error. buf:%v", buf))
	}
	return nil
}

//serialize the consensus payload
func (this *ConsensusPayload) ToArray() []byte {
	b := new(bytes.Buffer)
	err := this.Serialize(b)
	if err != nil {
		log.Errorf("consensus payload serialize error in ToArray(). payload:%v", this)
		return nil
	}
	return b.Bytes()
}

//return inventory type
func (this *ConsensusPayload) InventoryType() common.InventoryType {
	return common.CONSENSUS
}

func (this *ConsensusPayload) GetMessage() []byte {
	//TODO: GetMessage
	//return sig.GetHashData(cp)
	return []byte{}
}

func (this *ConsensusPayload) Type() common.InventoryType {

	//TODO:Temporary add for Interface signature.SignableData use.
	return common.CONSENSUS
}

//Serialize message payload
func (this *ConsensusPayload) Serialize(w io.Writer) error {
	err := this.SerializeUnsigned(w)
	if err != nil {
		return err
	}
	buf := keypair.SerializePublicKey(this.Owner)
	err = serialization.WriteVarBytes(w, buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write publickey error. publickey buf:%v", buf))
	}

	err = serialization.WriteVarBytes(w, this.Signature)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write Signature error. Signature:%v", this.Signature))
	}

	return nil
}

//Deserialize message payload
func (this *ConsensusPayload) Deserialize(r io.Reader) error {
	err := this.DeserializeUnsigned(r)
	if err != nil {
		return err
	}
	buf, err := serialization.ReadVarBytes(r)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read buf error")
	}
	this.Owner, err = keypair.DeserializePublicKey(buf)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "deserialize publickey error")
	}

	this.Signature, err = serialization.ReadVarBytes(r)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read Signature error")
	}

	return err
}

//Serialize message payload
func (this *ConsensusPayload) SerializeUnsigned(w io.Writer) error {
	err := serialization.WriteUint32(w, this.Version)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. version:%v", this.Version))
	}
	err = this.PrevHash.Serialize(w)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. PrevHash:%v", this.PrevHash))
	}
	err = serialization.WriteUint32(w, this.Height)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. Height:%v", this.Height))
	}
	err = serialization.WriteUint16(w, this.BookkeeperIndex)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. BookkeeperIndex:%v", this.BookkeeperIndex))
	}
	err = serialization.WriteUint32(w, this.Timestamp)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. Timestamp:%v", this.Timestamp))
	}
	err = serialization.WriteVarBytes(w, this.Data)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. Data:%v", this.Data))
	}
	return nil
}

//Deserialize message payload
func (this *ConsensusPayload) DeserializeUnsigned(r io.Reader) error {
	var err error
	this.Version, err = serialization.ReadUint32(r)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read version error")
	}

	preBlock := new(common.Uint256)
	err = preBlock.Deserialize(r)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read preBlock error")
	}
	this.PrevHash = *preBlock

	this.Height, err = serialization.ReadUint32(r)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read Height error")
	}

	this.BookkeeperIndex, err = serialization.ReadUint16(r)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read BookkeeperIndex error")
	}

	this.Timestamp, err = serialization.ReadUint32(r)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read Timestamp error")
	}

	this.Data, err = serialization.ReadVarBytes(r)
	if err != nil {

		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read Data error")
	}

	return nil
}
