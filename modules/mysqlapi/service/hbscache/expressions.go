package hbscache

import (
	"sync"

	"github.com/fwtpe/owl/common/model"
	db "github.com/fwtpe/owl/modules/mysqlapi/rdb/hbsdb"
)

type SafeExpressionCache struct {
	sync.RWMutex
	L []*model.NewExpression
}

var ExpressionCache = &SafeExpressionCache{}

func (this *SafeExpressionCache) Get() []*model.NewExpression {
	this.RLock()
	defer this.RUnlock()
	return this.L
}

func (this *SafeExpressionCache) Init() {
	es, err := db.QueryExpressions()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.L = es
}
