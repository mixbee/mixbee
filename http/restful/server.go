

package restful

import (
	"github.com/mixbee/mixbee/http/restful/restful"
)

func StartServer() {
	rt := restful.InitRestServer()
	go rt.Start()
}
