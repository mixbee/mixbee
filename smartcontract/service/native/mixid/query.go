
package mixid

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
)

func GetPublicKeyByID(srvc *native.NativeService) ([]byte, error) {
	args := bytes.NewBuffer(srvc.Input)
	// arg0: ID
	arg0, err := serialization.ReadVarBytes(args)
	if err != nil {
		return nil, errors.New("get public key failed: argument 0 error")
	}
	// arg1: key ID
	arg1, err := serialization.ReadUint32(args)
	if err != nil {
		return nil, errors.New("get public key failed: argument 1 error")
	}

	key, err := encodeID(arg0)
	if err != nil {
		return nil, fmt.Errorf("get public key failed: %s", err)
	}

	pk, err := getPk(srvc, key, arg1)
	if err != nil {
		return nil, fmt.Errorf("get public key failed: %s", err)
	} else if pk == nil {
		return nil, errors.New("get public key failed: not found")
	} else if pk.revoked {
		return nil, errors.New("get public key failed: revoked")
	}

	return pk.key, nil
}

func GetDDO(srvc *native.NativeService) ([]byte, error) {
	log.Debug("GetDDO")
	var0, err := GetPublicKeys(srvc)
	if err != nil {
		return nil, fmt.Errorf("get DDO error: %s", err)
	} else if var0 == nil {
		log.Debug("DDO: null")
		return nil, nil
	}
	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, var0)

	var1, err := GetAttributes(srvc)
	serialization.WriteVarBytes(&buf, var1)

	args := bytes.NewBuffer(srvc.Input)
	did, _ := serialization.ReadVarBytes(args)
	fmt.Printf("GetDDO arg0:%x\n",did)
	key, _ := encodeID(did)
	var2, err := getRecovery(srvc, key)
	serialization.WriteVarBytes(&buf, var2)

	res := buf.Bytes()
	log.Debug("DDO:", hex.EncodeToString(res))
	return res, nil
}

func GetPublicKeys(srvc *native.NativeService) ([]byte, error) {
	log.Debug("GetPublicKeys")
	args := bytes.NewBuffer(srvc.Input)
	did, err := serialization.ReadVarBytes(args)
	if err != nil {
		return nil, fmt.Errorf("get public keys error: invalid argument, %s", err)
	}
	if len(did) == 0 {
		return nil, errors.New("get public keys error: invalid ID")
	}
	key, err := encodeID(did)
	if err != nil {
		return nil, fmt.Errorf("get public keys error: %s", err)
	}
	key = append(key, FIELD_PK)
	list, err := getAllPk(srvc, key)
	if err != nil {
		return nil, fmt.Errorf("get public keys error: %s", err)
	} else if list == nil {
		return nil, nil
	}

	var res bytes.Buffer
	for i, v := range list {
		if v.revoked {
			continue
		}
		err = serialization.WriteUint32(&res, uint32(i+1))
		if err != nil {
			return nil, fmt.Errorf("get public keys error: %s", err)
		}
		err = serialization.WriteVarBytes(&res, v.key)
		if err != nil {
			return nil, fmt.Errorf("get public keys error: %s", err)
		}
	}

	return res.Bytes(), nil
}

func GetAttributes(srvc *native.NativeService) ([]byte, error) {
	log.Debug("GetAttributes")
	args := bytes.NewBuffer(srvc.Input)
	did, err := serialization.ReadVarBytes(args)
	if err != nil {
		return nil, fmt.Errorf("get public keys error: invalid argument", err)
	}
	if len(did) == 0 {
		return nil, errors.New("get attributes error: invalid ID")
	}
	key, err := encodeID(did)
	if err != nil {
		return nil, fmt.Errorf("get public keys error: %s", err)
	}
	res, err := getAllAttr(srvc, key)
	if err != nil {
		return nil, fmt.Errorf("get attributes error: %s", err)
	}

	return res, nil
}

func GetKeyState(srvc *native.NativeService) ([]byte, error) {
	log.Debug("GetKeyState")
	args := bytes.NewBuffer(srvc.Input)
	// arg0: ID
	arg0, err := serialization.ReadVarBytes(args)
	if err != nil {
		return nil, fmt.Errorf("get key state failed: argument 0 error, %s", err)
	}
	// arg1: public key ID
	arg1, err := utils.ReadVarUint(args)
	if err != nil {
		return nil, fmt.Errorf("get key state failed: argument 1 error, %s", err)
	}

	key, err := encodeID(arg0)
	if err != nil {
		return nil, fmt.Errorf("get key state failed: %s", err)
	}

	owner, err := getPk(srvc, key, uint32(arg1))
	if err != nil {
		return nil, fmt.Errorf("get key state failed: %s", err)
	} else if owner == nil {
		log.Debug("key state: not exist")
		return []byte("not exist"), nil
	}

	log.Debug("key state: ", owner.revoked)
	if owner.revoked {
		return []byte("revoked"), nil
	} else {
		return []byte("in use"), nil
	}
}
