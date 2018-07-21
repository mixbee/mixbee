

package payload

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mixbee/mixbee/common/serialization"
)

// DeployCode is an implementation of transaction payload for deploy smartcontract
type DeployCode struct {
	Code        []byte
	NeedStorage bool
	Name        string
	Version     string
	Author      string
	Email       string
	Description string
}

func (dc *DeployCode) Serialize(w io.Writer) error {
	var err error

	err = serialization.WriteVarBytes(w, dc.Code)
	if err != nil {
		return fmt.Errorf("DeployCode Code Serialize failed: %s", err)
	}

	err = serialization.WriteBool(w, dc.NeedStorage)
	if err != nil {
		return fmt.Errorf("DeployCode NeedStorage Serialize failed: %s", err)
	}

	err = serialization.WriteString(w, dc.Name)
	if err != nil {
		return fmt.Errorf("DeployCode Name Serialize failed: %s", err)
	}

	err = serialization.WriteString(w, dc.Version)
	if err != nil {
		return fmt.Errorf("DeployCode Version Serialize failed: %s", err)
	}

	err = serialization.WriteString(w, dc.Author)
	if err != nil {
		return fmt.Errorf("DeployCode Author Serialize failed: %s", err)
	}

	err = serialization.WriteString(w, dc.Email)
	if err != nil {
		return fmt.Errorf("DeployCode Email Serialize failed: %s", err)
	}

	err = serialization.WriteString(w, dc.Description)
	if err != nil {
		return fmt.Errorf("DeployCode Description Serialize failed: %s", err)
	}

	return nil
}

func (dc *DeployCode) Deserialize(r io.Reader) error {
	code, err := serialization.ReadVarBytes(r)
	if err != nil {
		return fmt.Errorf("DeployCode Code Deserialize failed: %s", err)
	}
	dc.Code = code

	dc.NeedStorage, err = serialization.ReadBool(r)
	if err != nil {
		return fmt.Errorf("DeployCode NeedStorage Deserialize failed: %s", err)
	}

	dc.Name, err = serialization.ReadString(r)
	if err != nil {
		return fmt.Errorf("DeployCode Name Deserialize failed: %s", err)
	}

	dc.Version, err = serialization.ReadString(r)
	if err != nil {
		return fmt.Errorf("DeployCode CodeVersion Deserialize failed: %s", err)
	}

	dc.Author, err = serialization.ReadString(r)
	if err != nil {
		return fmt.Errorf("DeployCode Author Deserialize failed: %s", err)
	}

	dc.Email, err = serialization.ReadString(r)
	if err != nil {
		return fmt.Errorf("DeployCode Email Deserialize failed: %s", err)
	}

	dc.Description, err = serialization.ReadString(r)
	if err != nil {
		return fmt.Errorf("DeployCode Description Deserialize failed: %s", err)
	}

	return nil
}

func (dc *DeployCode) ToArray() []byte {
	b := new(bytes.Buffer)
	dc.Serialize(b)
	return b.Bytes()
}
