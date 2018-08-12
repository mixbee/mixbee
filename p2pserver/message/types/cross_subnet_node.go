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

type CrossSubNetNodePayload struct {
	Host  string `json:"host"`
	NetId uint32 `json:"netId"`
}

func (this *CrossSubNetNodePayload) Serialize(w io.Writer) error {

	err := serialization.WriteString(w, this.Host)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write host error. host buf:%s", this.Host))
	}

	err = serialization.WriteUint32(w, this.NetId)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write netId error. netId:%d", this.NetId))
	}

	return nil
}

func (this *CrossSubNetNodePayload) Deserialize(r io.Reader) error {

	host, err := serialization.ReadString(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read host error")
	}
	this.Host = host

	this.NetId,err = serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "read netId error")
	}

	return err
}

//Serialize message payload
func (this *CrossSubNetNodePayload) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := this.Serialize(p)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. sub net node:%v", this))
	}
	return p.Bytes(), nil
}

func (this *CrossSubNetNodePayload) CmdType() string {
	return common.CROSSCHAIN_SUBNET_TYPE
}

//Deserialize message payload
func (this *CrossSubNetNodePayload) Deserialization(p []byte) error {
	log.Debug()
	buf := bytes.NewBuffer(p)
	err := this.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("deserialize sub net Node error. buf:%v", buf))
	}
	return nil
}
