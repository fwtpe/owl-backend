package receiver

import (
	"github.com/fwtpe/owl/modules/transfer/receiver/rpc"
)

func Start() {
	go rpc.StartRpc()
}
