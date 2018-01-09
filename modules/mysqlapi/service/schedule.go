package service

import (
	"fmt"
	"time"

	owlModel "github.com/fwtpe/owl-backend/common/model/owl"

	"github.com/fwtpe/owl-backend/common/utils"
	"github.com/fwtpe/owl-backend/modules/mysqlapi/model"
	"github.com/fwtpe/owl-backend/modules/mysqlapi/rdb"
)

// ScheduleService is designed to be a Execute function with namespace.
var ScheduleService = &scheduleService{
	Execute: ScheduleExecutor,
}

type ScheduleCallback func() error

type scheduleService struct {
	Execute func(*model.Schedule, ScheduleCallback) (*model.OwlScheduleLog, error)
}

func ScheduleExecutor(schedule *model.Schedule, callback ScheduleCallback) (*model.OwlScheduleLog, error) {
	scheduleLog, err := rdb.AcquireLock(schedule, time.Now())
	if err != nil {
		return nil, err
	}

	var callbackHandler = func() {
		var err error = nil

		defer func() {
			msg := ""

			p := recover()
			if p != nil {
				msg = fmt.Sprintf("Panic from scheduled callback: %v", p)
			} else if err != nil {
				msg = fmt.Sprintf("Error from scheduled callback: %v", err)
			}

			status := owlModel.JobDone
			if msg != "" {
				status = owlModel.JobFailed
				logger.Warnf("Execute task: [%v] has error: %s", schedule, msg)
			}

			rdb.FreeLock(scheduleLog, model.TaskStatus(status), msg, time.Now())
		}()

		err = callback()
	}

	go utils.BuildPanicCapture(
		callbackHandler,
		func(p interface{}) {
			logger.Errorf("During free lock of %s. Panic: %v", schedule, p)
		},
	)()

	return scheduleLog, nil
}
