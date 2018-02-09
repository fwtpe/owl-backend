package log

import (
	cgin "github.com/fwtpe/owl-backend/common/gin"
	"github.com/gin-gonic/gin"
)

type NamedLogger struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

func (l *NamedLogger) Bind(context *gin.Context) {
	cgin.BindJson(context, l)
}

type NamedLoggerList struct {
	Loggers []*NamedLogger `json:"loggers"`
}
