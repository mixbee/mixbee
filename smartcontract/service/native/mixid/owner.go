
package mixid

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/core/states"
	"github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
)

type owner struct {
	key     []byte
	revoked bool
}

func (this *owner) Serialize(w io.Writer) error {
	if err := serialization.WriteVarBytes(w, this.key); err != nil {
		return err
	}
	if err := serialization.WriteBool(w, this.revoked); err != nil {
		return err
	}
	return nil
}

func (this *owner) Deserialize(r io.Reader) error {
	v1, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	v2, err := serialization.ReadBool(r)
	if err != nil {
		return err
	}
	this.key = v1
	this.revoked = v2
	return nil
}

func getAllPk(srvc *native.NativeService, key []byte) ([]*owner, error) {

	val, err := utils.GetStorageItem(srvc, key)
	if err != nil {
		return nil, fmt.Errorf("get storage error, %s", err)
	}
	if val == nil {
		return nil, nil
	}
	buf := bytes.NewBuffer(val.Value)
	owners := make([]*owner, 0)
	for buf.Len() > 0 {
		var t = new(owner)
		err = t.Deserialize(buf)
		if err != nil {
			return nil, fmt.Errorf("deserialize owners error, %s", err)
		}
		owners = append(owners, t)
	}
	return owners, nil
}

func putAllPk(srvc *native.NativeService, key []byte, val []*owner) error {
	var buf bytes.Buffer
	for _, i := range val {
		err := i.Serialize(&buf)
		if err != nil {
			return fmt.Errorf("serialize owner error, %s", err)
		}
	}
	var v states.StorageItem
	v.Value = buf.Bytes()
	srvc.CloneCache.Add(common.ST_STORAGE, key, &v)
	return nil
}

func insertPk(srvc *native.NativeService, encID, pk []byte) (uint32, error) {
	key := append(encID, FIELD_PK)
	owners, err := getAllPk(srvc, key)
	if err != nil {
		owners = make([]*owner, 0)
	}
	size := len(owners)
	if size >= 0xFFFFFFFF {
		//FIXME currently the limit is for all the keys, including the
		//      revoked ones.
		return 0, errors.New("reach the max limit, cannot add more keys")
	}
	owners = append(owners, &owner{pk, false})
	err = putAllPk(srvc, key, owners)
	if err != nil {
		return 0, err
	}
	return uint32(size + 1), nil
}

func getPk(srvc *native.NativeService, encID []byte, index uint32) (*owner, error) {
	key := append(encID, FIELD_PK)
	owners, err := getAllPk(srvc, key)
	if err != nil {
		return nil, err
	}
	if index < 1 || index > uint32(len(owners)) {
		return nil, errors.New("invalid key index")
	}
	return owners[index-1], nil
}

func findPk(srvc *native.NativeService, encID, pub []byte) (uint32, error) {
	key := append(encID, FIELD_PK)
	owners, err := getAllPk(srvc, key)
	if err != nil {
		return 0, err
	}
	for i, v := range owners {
		if bytes.Equal(pub, v.key) {
			return uint32(i + 1), nil
		}
	}
	return 0, nil
}

func revokePk(srvc *native.NativeService, encID, pub []byte) (uint32, error) {
	key := append(encID, FIELD_PK)
	owners, err := getAllPk(srvc, key)
	if err != nil {
		return 0, err
	}
	var index uint32 = 0
	for i, v := range owners {
		if bytes.Equal(pub, v.key) {
			index = uint32(i + 1)
			if v.revoked {
				return index, errors.New("public key has already been revoked")
			}
			v.revoked = true
		}
	}
	if index == 0 {
		return 0, errors.New("revoke failed, public key not found")
	}
	err = putAllPk(srvc, key, owners)
	if err != nil {
		return 0, err
	}
	return index, nil
}

func isOwner(srvc *native.NativeService, encID, pub []byte) bool {
	kID, err := findPk(srvc, encID, pub)
	if err != nil {
		log.Debug(err)
		return false
	}
	if kID == 0 {
		return false
	}
	return true
}
