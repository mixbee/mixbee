

package types

import (
	"bytes"
	"fmt"

	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/p2pserver/common"
	"io"
	"github.com/mixbee/mixbee/common/serialization"
)

type CrossVerifyNodePayload struct {
	Node CrossChainVerifyNode
}
type CrossChainVerifyNode struct {
	PublicKey string `json:"publicKey"`
	Host      string `json:"host"`
	Type      uint32 `json:"type"` //1 register 2 delete
}

func (this *CrossChainVerifyNode) Serialize(w io.Writer) error {

	err := serialization.WriteString(w, this.PublicKey)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write publickey error. publickey buf:%s", this.PublicKey))
	}

	err = serialization.WriteString(w, this.Host)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write host error. Host:%s", this.Host))
	}

	err = serialization.WriteUint32(w, this.Type)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write type error. type:%d", this.Type))
	}

	return nil
}

func (this *CrossChainVerifyNode) Deserialize(r io.Reader) error {

	pbk, err := serialization.ReadString(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read public key error")
	}
	this.PublicKey = pbk

	this.Host, err = serialization.ReadString(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read Host error")
	}

	this.Type, err = serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read Type error")
	}

	return err
}

//Serialize message payload
func (this *CrossVerifyNodePayload) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := this.Node.Serialize(p)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. consensus:%v", this.Node))
	}
	return p.Bytes(), nil
}

func (this *CrossVerifyNodePayload) CmdType() string {
	return common.CROSSCHAIN_TYPE
}

//Deserialize message payload
func (this *CrossVerifyNodePayload) Deserialization(p []byte) error {
	log.Debug()
	buf := bytes.NewBuffer(p)
	err := this.Node.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("deserialize Node error. buf:%v", buf))
	}
	return nil
}
