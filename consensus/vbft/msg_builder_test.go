

package vbft

import (
	"fmt"
	"testing"

	"github.com/mixbee/mixbee/account"
)

func constructMsg() *blockProposalMsg {
	acc := account.NewAccount("SHA256withECDSA")
	if acc == nil {
		fmt.Println("GetDefaultAccount error: acc is nil")
		return nil
	}
	msg := constructProposalMsgTest(acc)
	return msg
}
func TestSerializeVbftMsg(t *testing.T) {
	msg := constructMsg()
	_, err := SerializeVbftMsg(msg)
	if err != nil {
		t.Errorf("TestSerializeVbftMsg failed :%v", err)
		return
	}
	t.Logf("TestSerializeVbftMsg succ")
}

func TestDeserializeVbftMsg(t *testing.T) {
	msg := constructMsg()
	data, err := SerializeVbftMsg(msg)
	if err != nil {
		t.Errorf("TestSerializeVbftMsg failed :%v", err)
		return
	}
	_, err = DeserializeVbftMsg(data)
	if err != nil {
		t.Errorf("DeserializeVbftMsg failed :%v", err)
		return
	}
	t.Logf("TestDeserializeVbftMsg succ")
}
