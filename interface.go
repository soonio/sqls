package sqls

import (
	"context"
)

type Lister[D any] interface {
	Querier[D]
	Counter[D]
	HasRecord() bool
}

type Querier[D any] interface {
	Query(context.Context, D) error
}

type Counter[D any] interface {
	Count(context.Context, D) error
}

type Pager interface {
	Limit() int64
	Offset() int64
}

type Sorter interface {
	OrderBy(...func(string) string) string
}

type Filter interface {
	Condition(only ...bool) (string, []any)
}
