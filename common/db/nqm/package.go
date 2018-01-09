package nqm

import (
	f "github.com/fwtpe/owl/common/db/facade"
	log "github.com/fwtpe/owl/common/logruslog"
	tb "github.com/fwtpe/owl/common/textbuilder"
)

var DbFacade *f.DbFacade

var t = tb.Dsl

var logger = log.NewDefaultLogger("warn")
