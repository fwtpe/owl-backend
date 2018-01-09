package restful

import (
	"github.com/fwtpe/owl-backend/common/gin/mvc"
	"github.com/fwtpe/owl-backend/modules/mysqlapi/rdb"
)

func getAgentConfig(
	q *struct {
		Key string `mvc:"query[key]"`
	},
) mvc.OutputBody {
	if q.Key == "" {
		return mvc.NotFoundOutputBody
	}
	retBody := rdb.GetAgentConfig(q.Key)
	return mvc.JsonOutputOrNotFound(retBody)
}
