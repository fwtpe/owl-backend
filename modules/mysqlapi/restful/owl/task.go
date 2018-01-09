package owl

import (
	"time"

	"github.com/fwtpe/owl-backend/common/gin/mvc"
	owlDb "github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/owl"
)

func ClearLogsOfTasks(
	p *struct {
		ForDays int `mvc:"query[for_days] default[90]" validate:"min=1"`
	},
) mvc.OutputBody {
	t := time.Now().Add(time.Duration(-p.ForDays) * time.Duration(24) * time.Hour)
	affectedRows := owlDb.RemoveOldScheduleLogs(t)

	return mvc.JsonOutputBody(
		map[string]interface{}{
			"before_time":   t.Unix(),
			"affected_rows": affectedRows,
		},
	)
}
