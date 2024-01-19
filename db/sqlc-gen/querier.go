// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"
)

type Querier interface {
	// :one
	SchemaDelete(ctx context.Context, id int32) (Schema, error)
	// :one
	SchemaGet(ctx context.Context, id int32) (Schema, error)
	SchemaList(ctx context.Context, arg SchemaListParams) ([]Schema, error)
	// :one
	SchemaNew(ctx context.Context, arg SchemaNewParams) (Schema, error)
	// :one
	SchemaUpdate(ctx context.Context, arg SchemaUpdateParams) (Schema, error)
}

var _ Querier = (*Queries)(nil)
