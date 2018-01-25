package database

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	cdb "github.com/fwtpe/owl-backend/common/db"
	"github.com/fwtpe/owl-backend/common/db/facade"
	nqmDb "github.com/fwtpe/owl-backend/common/db/nqm"
	owlDb "github.com/fwtpe/owl-backend/common/db/owl"
	oHttp "github.com/fwtpe/owl-backend/common/http"
	owlSrv "github.com/fwtpe/owl-backend/common/service/owl"

	"github.com/fwtpe/owl-backend/modules/query/g"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var PortalDbFacade *facade.DbFacade

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
	PortalDbFacade = &facade.DbFacade{}
	err = PortalDbFacade.Open(
		&cdb.DbConfig{
			Dsn:     conf.Db.Addr,
			MaxIdle: conf.Db.Idle,
		},
	)

	if err != nil {
		log.Printf("%v\n", err)
	}

	owlDb.DbFacade = PortalDbFacade
	nqmDb.DbFacade = PortalDbFacade
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
