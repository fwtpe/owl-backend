package log

import (
	mvc "github.com/fwtpe/owl-backend/common/gin/mvc"
	cModel "github.com/fwtpe/owl-backend/common/model"
	"github.com/gin-gonic/gin"
)

func RestLogger(prefixGroup gin.IRouter, b mvc.MvcBuilder) {
	h := b.BuildHandler
	prefixGroup.GET("/v1/list-all", h(RestListAllV1))
	prefixGroup.POST("/v1/set-level", h(RestSetLevelV1))
}

func RestListAllV1() mvc.OutputBody {
	loggers := ListAll()
	reply := cModel.NamedLoggerList{make([]*cModel.NamedLogger, 0, len(loggers))}
	for _, entry := range loggers {
		reply.Loggers = append(reply.Loggers, &cModel.NamedLogger{
			Name:  entry.Name,
			Level: entry.Logger.Level.String(),
		})
	}
	return mvc.JsonOutputBody(reply)
}

func RestSetLevelV1(
	setData cModel.NamedLogger,
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
