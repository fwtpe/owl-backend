package cmdb

import (
	"github.com/fwtpe/owl-backend/modules/mysqlapi/model"
	bossRdb "github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/boss"
	cmdbRdb "github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/cmdb"
	"github.com/fwtpe/owl-backend/modules/mysqlapi/service"
	bossSrv "github.com/fwtpe/owl-backend/modules/mysqlapi/service/boss"
)

func SyncDataFromBoss() (*model.OwlScheduleLog, error) {
	return service.ScheduleService.Execute(
		model.NewSchedule("import.imdb", 300),
		func() error {
			sourceData := bossSrv.Boss2cmdb(bossRdb.GetSyncData())

			cmdbRdb.SyncForHosts(sourceData)

			return nil
		},
	)
}
