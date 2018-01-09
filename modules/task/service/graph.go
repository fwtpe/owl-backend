package service

import (
	graphSrv "github.com/fwtpe/owl/common/service/graph"

	"github.com/fwtpe/owl/modules/task/database"
	"github.com/fwtpe/owl/modules/task/proc"
)

func VacuumGraphIndex(beforeDays int) *graphSrv.ResultOfVacuumIndex {
	defer func() {
		proc.IndexDeleteCnt.Incr()
	}()

	/**
	 * Puts "0" counter for deletion of graph index
	 */
	proc.IndexDeleteCnt.PutOther("deleteCntEndpoint", 0)
	proc.IndexDeleteCnt.PutOther("deleteCntTagEndpoint", 0)
	proc.IndexDeleteCnt.PutOther("deleteCntEndpointCounter", 0)
	// :~)

	result := database.GraphService.VacuumIndex(
		&graphSrv.VacuumIndexConfig{
			BeforeDays: beforeDays,
		},
	)

	/**
	 * Puts actual counter for deletion of graph index
	 */
	affectedRows := result.AffectedRows
	proc.IndexDeleteCnt.PutOther("deleteCntEndpoint", affectedRows.Endpoints)
	proc.IndexDeleteCnt.PutOther("deleteCntTagEndpoint", affectedRows.Tags)
	proc.IndexDeleteCnt.PutOther("deleteCntEndpointCounter", affectedRows.Counters)
	// :~)

	return result
}
