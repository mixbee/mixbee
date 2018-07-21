
package handlers

import (
	"encoding/hex"
	"encoding/json"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	"testing"
)

func TestSigData(t *testing.T) {
	rawData := []byte("HelloWorld")
	rawReq := &SigDataReq{
		RawData: hex.EncodeToString(rawData),
	}
	data, err := json.Marshal(rawReq)
	if err != nil {
		t.Errorf("json.Marshal SigDataReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "sigdata",
		Params: data,
	}
	resp := &clisvrcom.CliRpcResponse{}
	SigData(req, resp)
	if resp.ErrorCode != 0 {
		t.Errorf("SigData failed. ErrorCode:%d", resp.ErrorCode)
		return
	}
}
