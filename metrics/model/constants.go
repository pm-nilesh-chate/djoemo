package model

import "context"

const (
	StatusSuccess = "success"
	StatusFailure = "failure"

	OpCommit = "commit"
	OpUpdate = "update"
	OpRead   = "read"
	OpDelete = "delete"
)

type djoemoMetricsContextKey string

const (
	ContextKeySource djoemoMetricsContextKey = "source"
)

// WithSourceLabel is a label to tag buisness logic as default metrics are aggregated for CURD operations to reduce cardinality
func WithSourceLabel(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, ContextKeySource, name)
}
