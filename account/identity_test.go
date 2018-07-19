package account

import (
	"testing"
	"encoding/hex"
)

var testId = "TPFtz8vXZoGrmTkyRqg5ZXELdePCB2hDm5"


func TestCreateID(t *testing.T)  {
	nonce, _ := hex.DecodeString("4c6b58adc6b8c6774eee0eb07dac4e198df87aae28f8932db3982edf3ff026e4")
	id, err := CreateID(nonce)
	if err!=nil {
		t.Fatal(err)
	}
	t.Log("result ID:", id)


}

func TestGenerateID(t *testing.T)  {
	id,err := GenerateID()
	if err!=nil {
		t.Fatal(err)
	}
	t.Log("result ID:", id)
}

func TestVerifyID(t *testing.T)  {
	if !VerifyID(testId) {
		t.Error("error: failed")
	}

}

func TestNewIdentity(t *testing.T)  {

}

