package sqls

import (
	"reflect"
	"strings"
)

type R struct {
	k []string // 数据表字段的键
	v []any    // 数据表字段的值列表
}

func (r *R) Scope(f func(*R)) *R {
	f(r)
	return r
}

// When 当条件成立时才添加
func (r *R) When(do bool, k string, value ...any) *R {
	if do {
		r.Add(k, value...)
	}
	return r
}

// Pointer 判断指针中是否有值
func (r *R) Pointer(sql string, value any) *R {
	return r.When(!reflect.ValueOf(value).IsNil(), sql, value)
}

func (r *R) In(column string, items ...any) *R {
	if len(items) == 0 {
		return r
	}

	if len(items) == 1 {
		r.Add(column+" = ?", items[0])
		return r
	}

	var sb strings.Builder
	sb.WriteString(column)
	sb.WriteString(" in (?")
	for i := 1; i < len(items); i++ {
		sb.WriteString(",?")
	}
	sb.WriteString(")")

	r.Add(sb.String(), items...)

	return r

}

func (r *R) Add(k string, v ...any) *R {
	r.k = append(r.k, k)
	r.v = append(r.v, v...)
	return r
}

func (r *R) SQL(where ...bool) (string, []any) {
	return r.Build(where...), r.v
}

func (r *R) Build(where ...bool) string {
	where = append(where, false)
	var sb strings.Builder
	if len(r.k) > 0 {
		if where[0] {
			sb.WriteString(" where ")
		}
		sb.WriteString(r.k[0])
	}

	for i := 1; i < len(r.k); i++ {
		sb.WriteString(" and ")
		sb.WriteString(r.k[i])
	}
	return sb.String()
}
