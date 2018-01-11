package database

import (
	"github.com/jinzhu/gorm"

	cdb "github.com/fwtpe/owl-backend/common/db"
	"github.com/fwtpe/owl-backend/common/db/facade"
	nqmDb "github.com/fwtpe/owl-backend/common/db/nqm"
	owlDb "github.com/fwtpe/owl-backend/common/db/owl"
	oHttp "github.com/fwtpe/owl-backend/common/http"
	owlSrv "github.com/fwtpe/owl-backend/common/service/owl"
	"github.com/fwtpe/owl-backend/common/logruslog"

	"github.com/fwtpe/owl-backend/modules/query/g"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var PortalDbFacade *facade.DbFacade
var BossDbFacade *facade.DbFacade

var logger = logruslog.NewDefaultLogger("INFO")

var (
	db  *gorm.DB
	err error
)

func DBConn() *gorm.DB {
	return db
}

func Init() {
	conf := g.Config()

	/**
	 * Use Db Facade to initialize related service
	 */
	PortalDbFacade = openDbFacade(
		&cdb.DbConfig{
			Dsn:     conf.Db.Addr,
			MaxIdle: conf.Db.Idle,
		},
		"portal",
	)
	owlDb.DbFacade = PortalDbFacade
	nqmDb.DbFacade = PortalDbFacade

	BossDbFacade = openDbFacade(
		&cdb.DbConfig{
			Dsn:     conf.BossDB.Addr,
			MaxIdle: conf.BossDB.Idle,
		},
		"boss",
	)
	// :~)

	db = PortalDbFacade.GormDb
}

var (
	QueryObjectService owlSrv.QueryService
)

func InitMySqlApi(config *oHttp.RestfulClientConfig) {
	QueryObjectService = owlSrv.NewQueryService(
		owlSrv.QueryServiceConfig{
			config,
		},
	)
}

func openDbFacade(config *cdb.DbConfig, name string) *facade.DbFacade {
	newDbFacade := &facade.DbFacade{}

	if err := newDbFacade.Open(config); err != nil {
		logger.Errorf("Cannot open MySql to [%s]: %v", name, err)
		return nil
	}

	return newDbFacade
}
