package restful

import (
	"github.com/fwtpe/owl/common/gin/mvc"
	apiModel "github.com/fwtpe/owl/common/model/mysqlapi"
	"github.com/fwtpe/owl/modules/mysqlapi/rdb"
	"github.com/fwtpe/owl/modules/mysqlapi/service"
)

func health() mvc.OutputBody {
	portalRdbDiag := rdb.GlobalDbHolder.Diagnose(rdb.DB_PORTAL)
	graphRdbDiag := rdb.GlobalDbHolder.Diagnose(rdb.DB_GRAPH)
	bossRdbDiag := rdb.GlobalDbHolder.Diagnose(rdb.DB_BOSS)

	health := &apiModel.HealthView{
		Rdb: &apiModel.AllRdbHealth{
			Dsn:             portalRdbDiag.Dsn,
			OpenConnections: portalRdbDiag.OpenConnections,
			PingResult:      portalRdbDiag.PingResult,
			PingMessage:     portalRdbDiag.PingMessage,
			Portal:          portalRdbDiag,
			Graph:           graphRdbDiag,
			Boss:            bossRdbDiag,
		},
		Http: &apiModel.Http{
			Listening: GinConfig.GetAddress(),
		},
		Nqm: &apiModel.Nqm{
			Heartbeat: &apiModel.Heartbeat{
				Count: service.NqmQueue.ConsumedCount(),
			},
		},
	}

	return mvc.JsonOutputBody(health)
}
