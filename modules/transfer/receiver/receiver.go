package receiver

import (
	"github.com/fwtpe/owl-backend/modules/transfer/receiver/rpc"
)

func Start() {
	go rpc.StartRpc()
}
