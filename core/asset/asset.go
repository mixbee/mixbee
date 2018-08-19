package asset

import (
	"io"
	"github.com/mixbee/mixbee/common/serialization"
	. "github.com/mixbee/mixbee/errors"
	"bytes"
)

type AssetType byte

const (
	Currency AssetType = 0x00
	Share    AssetType = 0x01
	Invoice  AssetType = 0x10
	Token    AssetType = 0x11
)

const (
	MaxPrecision = 8
	MinPrecision = 0
)

type AssetRecordType byte


const (
	UTXO    AssetRecordType = 0x00
	Balance AssetRecordType = 0x01
)

type Asset struct {
	Name        string
	Description string
	Precision   byte
	AssetType   AssetType
	RecordType  AssetRecordType
}

func (a *Asset)Serialize(w io.Writer) error {
	err := serialization.WriteString(w, a.Name)
	if err!= nil {
		return  NewDetailErr(err, ErrNoCode, "[Asset], Name serialize failed.")
	}

	err = serialization.WriteString(w, a.Description)
	if err != nil {
		return  NewDetailErr(err, ErrNoCode, "[Asset], Description serialize failed.")
	}

	err = serialization.WriteByte(w, a.Precision)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], Precision serialize failed.")
	}

	err = serialization.WriteUint8(w, byte(a.AssetType))
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], AssetType serialize failed.")
	}

	err = serialization.WriteUint8(w, byte(a.RecordType))
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], RecordType serialize failed.")
	}
	return nil
}

func (a *Asset)Deserialize(r io.Reader) error  {
	name, err := serialization.ReadString(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode,"[Asset], Name deserialize failed.")
	}
	a.Name = name

	description,err := serialization.ReadString(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode,"[Asset], Description deserialize failed.")
	}
	a.Description = description

	precision,err := serialization.ReadByte(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode,"[Asset], Precision deserialize failed.")
	}
	a.Precision = precision

	val,err := serialization.ReadBytes(r, 1)
	if err != nil {
		return NewDetailErr(err, ErrNoCode,"[Asset], AssetType deserialize failed.")
	}
	a.AssetType = AssetType(val[0])

	val,err = serialization.ReadBytes(r, 1)
	if err != nil {
		return NewDetailErr(err, ErrNoCode,"[Asset], RecordType deserialize failed.")
	}
	a.RecordType = AssetRecordType(val[0])
	return nil
}

func (a *Asset)ToArray() ([]byte)  {
	b := new(bytes.Buffer)
	a.Serialize(b)
	return b.Bytes()
}


