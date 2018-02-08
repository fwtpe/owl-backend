package log

import (
	mvc "github.com/fwtpe/owl-backend/common/gin/mvc"
	"github.com/fwtpe/owl-backend/common/model"
	"github.com/gin-gonic/gin"
)

func RestLogger(prefixGroup gin.IRouter, h mvc.BuildHandlerFunc) {
	prefixGroup.GET("/v1/list-all", h(restListAllV1))
	prefixGroup.POST("/v1/set-level", h(restSetLevelV1))
}

func restListAllV1() mvc.OutputBody {
	loggers := ListAll()
	reply := model.NamedLoggerList{make([]*model.NamedLogger, 0, len(loggers))}
	for _, entry := range loggers {
		reply.Loggers = append(reply.Loggers, &model.NamedLogger{
			Name:  entry.Name,
			Level: entry.Logger.Level.String(),
		})
	}
	return mvc.JsonOutputBody(reply)
}

func restSetLevelV1(
	setData model.NamedLogger,
	q *struct {
		Tree bool `mvc:"query[tree]"`
	},
) mvc.OutputBody {
	var count int
	level, err := parseLevel(setData.Level)
	if err != nil {
		panic(err.Error())
	}

	if q.Tree {
		count = SetLevelToTree(setData.Name, level)
	} else {
		if ret := SetLevel(setData.Name, level); ret {
			count = 1
		}
	}

	return mvc.JsonOutputBody(gin.H{"affected_loggers": count})
}
