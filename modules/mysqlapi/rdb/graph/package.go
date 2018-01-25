package graph

import (
	f "github.com/fwtpe/owl-backend/common/db/facade"
	log "github.com/fwtpe/owl-backend/common/logruslog"
)

var DbFacade *f.DbFacade
var logger = log.NewDefaultLogger("warn")
