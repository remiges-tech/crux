// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"
)

type Querier interface {
	// :one
	SchemaDelete(ctx context.Context, id int32) (int32, error)
	SchemaGet(ctx context.Context, arg SchemaGetParams) ([]Schema, error)
	SchemaList(ctx context.Context) ([]SchemaListRow, error)
	SchemaListByApp(ctx context.Context, app string) ([]SchemaListByAppRow, error)
	SchemaListByAppAndClass(ctx context.Context, arg SchemaListByAppAndClassParams) ([]SchemaListByAppAndClassRow, error)
	SchemaListByAppAndSlice(ctx context.Context, arg SchemaListByAppAndSliceParams) ([]SchemaListByAppAndSliceRow, error)
	SchemaListByClass(ctx context.Context, class string) ([]SchemaListByClassRow, error)
	SchemaListByClassAndSlice(ctx context.Context, arg SchemaListByClassAndSliceParams) ([]SchemaListByClassAndSliceRow, error)
	SchemaListBySlice(ctx context.Context, slice int32) ([]SchemaListBySliceRow, error)
	// :one
	SchemaNew(ctx context.Context, arg SchemaNewParams) (int32, error)
	// :one
	SchemaUpdate(ctx context.Context, arg SchemaUpdateParams) (int32, error)
}

var _ Querier = (*Queries)(nil)
