package restful

import (
	mvc "github.com/fwtpe/owl/common/gin/mvc"
	commonModel "github.com/fwtpe/owl/common/model"
	"github.com/fwtpe/owl/modules/mysqlapi/rdb"
)

func listHosts(
	paging *struct {
		Page *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[id#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	agents, resultPaging := rdb.ListHosts(*paging.Page)

	return resultPaging, mvc.JsonOutputBody(agents)
}

func listHostgroups(
	paging *struct {
		Page *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[id#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	hostgroups, resultPaging := rdb.ListHostgroups(*paging.Page)

	return resultPaging, mvc.JsonOutputBody(hostgroups)
}
