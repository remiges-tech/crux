package markdone

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	sqlc "github.com/remiges-tech/crux/db/sqlc-gen"
)

func deleteWFInstance(queries *sqlc.Queries, entity Markdone_t) error {
	params := sqlc.DeleteWFInstancesParams{
		Slice: entity.Entity.Slice,
		App:   entity.Entity.App,
		ID:    entity.InstanceID,
	}
	return queries.DeleteWFInstances(context.Background(), params)
}

func GetWFInstanceCountForEntity(queries *sqlc.Queries, entity Markdone_t, workflowname string) (int64, error) {
	params := sqlc.GetWFInstanceCountsParams{
		Slice:    entity.Entity.Slice,
		App:      entity.Entity.App,
		Workflow: workflowname,
		ID:       entity.InstanceID,
	}
	count, err := queries.GetWFInstanceCounts(context.Background(), params)
	if err != nil {
		log.Printf("Error running GetWFInstanceCounts: %v", err)
	}
	return count, err

}
func UpdateWFInstanceStep(queries *sqlc.Queries, entity Markdone_t, step string, workflowname string) error {

	params := sqlc.UpdateWFInstanceStepParams{
		Slice:    entity.Entity.Slice,
		App:      entity.Entity.App,
		ID:       int32(entity.InstanceID),
		Step:     step,
		Workflow: workflowname,
	}

	return queries.UpdateWFInstanceStep(context.Background(), params)

}
func UpdateWFInstanceDoneAt(queries *sqlc.Queries, entity Markdone_t, t time.Time, wf string) error {

	// id := strconv.Itoa(int(entity.InstanceID))
	params := sqlc.UpdateWFInstanceDoneatParams{
		Doneat:   pgtype.Timestamp{Time: t, Valid: true},
		ID:       entity.InstanceID,
		Slice:    entity.Entity.Slice,
		App:      entity.Entity.App,
		Workflow: wf,
	}

	return queries.UpdateWFInstanceDoneat(context.Background(), params)

}

func getWFInstanceList(queries *sqlc.Queries, entity Markdone_t, wf string) ([]sqlc.Wfinstance, error) {

	parent := &pgtype.Int4{} // Ensure parent is of type pgx/v5/pgtype.Int4
	parent.Int32 = entity.InstanceID
	params := sqlc.GetWFInstanceListForMarkDoneParams{
		ID:       entity.InstanceID,
		Slice:    entity.Entity.Slice,
		App:      entity.Entity.App,
		Workflow: wf,
	}
	list, err := queries.GetWFInstanceListForMarkDone(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func getCurrentWFINstance(queries *sqlc.Queries, entity Markdone_t, wf string) (sqlc.Wfinstance, error) {

	id := strconv.Itoa(int(entity.InstanceID))

	params := sqlc.GetWFInstanceCurrentParams{
		Entityid: id,
		Slice:    entity.Entity.Slice,
		App:      entity.Entity.App,
		Workflow: wf,
	}
	return queries.GetWFInstanceCurrent(context.Background(), params)

}

func GetSubFLow(queries *sqlc.Queries, step string) (sqlc.GetWorkflowNameForStepRow, error) {
	return queries.GetWorkflowNameForStep(context.Background(), step)
}
