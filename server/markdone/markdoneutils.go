package markdone

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	sqlc "github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
)

func deleteWFInstance(queries *sqlc.Queries, instanceID int32, entity crux.Entity) error {
	params := sqlc.DeleteWFInstancesParams{
		Slice: entity.Slice,
		App:   entity.App,
		ID:    instanceID,
	}
	return queries.DeleteWFInstances(context.Background(), params)
}

func GetWFInstanceCountForEntity(queries *sqlc.Queries, instanceID int32, entity crux.Entity, workflowname string) (int64, error) {
	params := sqlc.GetWFInstanceCountsParams{
		Slice:    entity.Slice,
		App:      entity.App,
		Workflow: workflowname,
		ID:       instanceID,
	}
	count, err := queries.GetWFInstanceCounts(context.Background(), params)
	if err != nil {
		log.Printf("Error running GetWFInstanceCounts: %v", err)
	}
	return count, err

}
func UpdateWFInstanceStep(queries *sqlc.Queries, instanceID int32, entity crux.Entity, step string, workflowname string) error {

	params := sqlc.UpdateWFInstanceStepParams{
		Slice:    entity.Slice,
		App:      entity.App,
		ID:       int32(instanceID),
		Step:     step,
		Workflow: workflowname,
	}

	return queries.UpdateWFInstanceStep(context.Background(), params)

}
func UpdateWFInstanceDoneAt(queries *sqlc.Queries, instanceID int32, entity crux.Entity, t time.Time, wf string) error {

	// id := strconv.Itoa(int(instanceID))
	params := sqlc.UpdateWFInstanceDoneatParams{
		Doneat:   pgtype.Timestamp{Time: t, Valid: true},
		ID:       instanceID,
		Slice:    entity.Slice,
		App:      entity.App,
		Workflow: wf,
	}

	return queries.UpdateWFInstanceDoneat(context.Background(), params)

}

func getWFInstanceList(queries *sqlc.Queries, instanceID int32, entity crux.Entity, wf string) ([]sqlc.Wfinstance, error) {

	parent := &pgtype.Int4{} // Ensure parent is of type pgx/v5/pgtype.Int4
	parent.Int32 = instanceID
	params := sqlc.GetWFInstanceListForMarkDoneParams{
		ID:       instanceID,
		Slice:    entity.Slice,
		App:      entity.App,
		Workflow: wf,
	}
	list, err := queries.GetWFInstanceListForMarkDone(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func getCurrentWFINstance(queries *sqlc.Queries, instanceID int32, entity crux.Entity, wf string) (sqlc.Wfinstance, error) {

	id := strconv.Itoa(int(instanceID))

	params := sqlc.GetWFInstanceCurrentParams{
		Entityid: id,
		Slice:    entity.Slice,
		App:      entity.App,
		Workflow: wf,
	}
	return queries.GetWFInstanceCurrent(context.Background(), params)

}

func GetSubFLow(queries *sqlc.Queries, step string) (sqlc.GetWorkflowNameForStepRow, error) {
	return queries.GetWorkflowNameForStep(context.Background(), step)
}
