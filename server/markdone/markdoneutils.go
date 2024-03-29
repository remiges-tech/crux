package markdone

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	sqlc "github.com/remiges-tech/crux/db/sqlc-gen"
)

var queryDbq *sqlc.Queries

func deleteWFInstance(entity Markdone_t) error {
	sliceInt, err := strconv.Atoi(entity.Entity.Slice)
	if err != nil {
		log.Fatal("Failed to convert string to int32:", err)
	}
	id := strconv.Itoa(int(entity.Id))
	params := sqlc.DeleteWFInstancesParams{
		Slice:    int32(sliceInt),
		App:      entity.Entity.App,
		Entityid: id,
	}
	return queryDbq.DeleteWFInstances(context.Background(), params)
}

func GetWFInstanceCountForEntity(queries *sqlc.Queries, entity Markdone_t, workflowname string) (int64, error) {
	queryDbq = queries
	sliceInt, err := strconv.Atoi(entity.Entity.Slice)
	if err != nil {
		log.Fatal("Failed to convert string to int32:", err)
		return -1, err
	}
	params := sqlc.GetWFInstanceCountsParams{
		Slice:    int32(sliceInt),
		App:      entity.Entity.App,
		Workflow: workflowname,
		ID:       entity.Id,
	}
	count, err := queryDbq.GetWFInstanceCounts(context.Background(), params)
	if err != nil {
		log.Printf("Error running GetWFInstanceCounts: %v", err)
	}
	return count, err

}
func UpdateWFInstanceStep(entity Markdone_t, step string) error {

	sliceInt, err := strconv.Atoi(entity.Entity.Slice)
	if err != nil {
		log.Fatal("Failed to convert string to int32:", err)
		return err
	}
	id := strconv.Itoa(int(entity.Id))
	params := sqlc.UpdateWFInstanceStepParams{
		Slice:    int32(sliceInt),
		App:      entity.Entity.App,
		Entityid: id,
		Step:     step,
	}

	return queryDbq.UpdateWFInstanceStep(context.Background(), params)

}
func UpdateWFInstanceDoneAt(entity Markdone_t, t time.Time, wf string) error {

	sliceInt, err := strconv.Atoi(entity.Entity.Slice)
	if err != nil {
		log.Fatal("Failed to convert string to int32:", err)
		return err
	}
	id := strconv.Itoa(int(entity.Id))
	params := sqlc.UpdateWFInstanceDoneatParams{
		Doneat:   pgtype.Timestamp{Time: t},
		Entityid: id,
		Slice:    int32(sliceInt),
		App:      entity.Entity.App,
		Workflow: wf,
	}

	return queryDbq.UpdateWFInstanceDoneat(context.Background(), params)

}

func getWFInstanceList(entity Markdone_t, wf string) ([]sqlc.Wfinstance, error) {

	sliceInt, err := strconv.Atoi(entity.Entity.Slice)
	if err != nil {
		log.Fatal("Failed to convert string to int32:", err)
		return nil, err
	}
	id := strconv.Itoa(int(entity.Id))
	parent := &pgtype.Int4{} // Ensure parent is of type pgx/v5/pgtype.Int4
	parent.Int32 = entity.Id
	params := sqlc.GetWFInstanceListParams{
		Entityid: pgtype.Text{String: id, Valid: true},
		Slice:    pgtype.Int4{Int32: int32(sliceInt), Valid: true},
		App:      pgtype.Text{String: entity.Entity.App, Valid: true},
		Workflow: pgtype.Text{String: wf, Valid: true},
		Parent:   *parent,
	}
	return queryDbq.GetWFInstanceList(context.Background(), params)
}
func getCurrentWFINstance(entity Markdone_t, wf string) (sqlc.Wfinstance, error) {
	sliceInt, err := strconv.Atoi(entity.Entity.Slice)
	if err != nil {
		log.Fatal("Failed to convert string to int32:", err)
		return sqlc.Wfinstance{}, err
	}
	id := strconv.Itoa(int(entity.Id))

	params := sqlc.GetWFInstanceCurrentParams{
		Entityid: id,
		Slice:    int32(sliceInt),
		App:      entity.Entity.App,
		Workflow: wf,
	}
	return queryDbq.GetWFInstanceCurrent(context.Background(), params)

}

func GetSubFLow(step string) (sqlc.GetWorkflowNameForStepRow, error) {
	return queryDbq.GetWorkflowNameForStep(context.Background(), step)
}
