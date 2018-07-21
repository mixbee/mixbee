

package localrpc

import (
	"net/http"
	"strconv"

	"fmt"
	cfg "github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/http/base/rpc"
)

const (
	LOCAL_HOST string = "127.0.0.1"
	LOCAL_DIR  string = "/local"
)

func StartLocalServer() error {
	log.Debug()
	http.HandleFunc(LOCAL_DIR, rpc.Handle)

	rpc.HandleFunc("getneighbor", rpc.GetNeighbor)
	rpc.HandleFunc("getnodestate", rpc.GetNodeState)
	rpc.HandleFunc("startconsensus", rpc.StartConsensus)
	rpc.HandleFunc("stopconsensus", rpc.StopConsensus)
	rpc.HandleFunc("setdebuginfo", rpc.SetDebugInfo)

	// TODO: only listen to local host
	err := http.ListenAndServe(":"+strconv.Itoa(int(cfg.DefConfig.Rpc.HttpLocalPort)), nil)
	if err != nil {
		return fmt.Errorf("ListenAndServe error:%s", err)
	}
	return nil
}
