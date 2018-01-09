package boss

import (
	f "github.com/fwtpe/owl/common/db/facade"
	log "github.com/fwtpe/owl/common/logruslog"
)

var DbFacade *f.DbFacade
var logger = log.NewDefaultLogger("warn")
