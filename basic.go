package sqls

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Datetime struct {
	time.Time
}

func (n *Datetime) Value() (driver.Value, error) {
	return n.Time.Format("2006-01-02 15:04:05"), nil
}

func (n *Datetime) Scan(value any) error {
	if value == nil {
		return nil
	}
	n.Time = value.(time.Time)
	return nil
}

func (n *Datetime) MarshalJSON() ([]byte, error) {
	if n.Time.IsZero() {
		return []byte(`""`), nil
	}
	b := make([]byte, 0, 21)
	b = append(b, '"')
	b = n.Time.AppendFormat(b, "2006-01-02 15:04:05")
	b = append(b, '"')
	return b, nil
}

func (n *Datetime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*n = Datetime{Time: now}
	return
}

func (n *Datetime) Format(layout string) string {
	return n.Time.Format(layout)
}

type Json[T any] struct {
	Data *T `json:"ignore"`
}

func (u *Json[T]) UnmarshalJSON(bytes []byte) error {
	return jsoniter.Unmarshal(bytes, &u.Data)
}

func (u *Json[T]) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(u.Data)
}

func (u *Json[T]) Scan(src any) error {
	if src == nil {
		return nil
	}
	if v, ok := src.([]byte); ok {
		return jsoniter.Unmarshal(v, u)
	}
	return fmt.Errorf("can not convert %v", src)
}

func (u *Json[T]) Value() (driver.Value, error) {
	return jsoniter.Marshal(u)
}

func Condition(r *R, p Pager, s Sorter, f ...func(string) string) (string, []any) {
	var sb strings.Builder
	c, args := r.SQL(true)
	sb.WriteString(c)

	if s != nil {
		ob := s.OrderBy(f...)
		if len(ob) > 0 {
			sb.WriteString(" ")
			sb.WriteString(ob)
		}
	}
	if p != nil {
		args = append(args, p.Limit(), p.Offset())
		sb.WriteString(" limit ? offset ?")
	}

	return sb.String(), args
}

type ListWithFilter[T any] struct {
	Filter Filter
	Total  int64
	List   []*T
}

func MustAffected(r sql.Result, err error) error {
	if err != nil {
		return err
	}
	if num, err := r.RowsAffected(); err != nil {
		return err
	} else {
		if num == 1 {
			return nil
		} else {
			return ErrNotAffected
		}
	}
}

func QueryLister[D any](ctx context.Context, d D, lister Lister[D]) error {
	var t = reflect.TypeOf(lister)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf("QueryLister requires a pointer")
	}

	var err = lister.Count(ctx, d)
	if err != nil {
		return err
	}
	if !lister.HasRecord() {
		return nil
	}
	err = lister.Query(ctx, d)
	if err != nil {
		return err
	}
	return nil
}
