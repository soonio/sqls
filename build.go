package sqls

import (
	"reflect"
	"slices"
	"strings"
)

// 结构体转换为数据
// 这里不转换为map而转换为切片，避免每次顺序不一样
func s2d(s any) ([]string, []any) {
	var keys []string
	var vals []any

	var t = reflect.TypeOf(s)
	var v = reflect.ValueOf(s)

	if t.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			keys = append(keys, key.String())
			vals = append(vals, v.MapIndex(key).Interface())
		}
		return keys, vals
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	fields := t.NumField()
	for i := 0; i < fields; i++ {
		key := t.Field(i).Name
		if t.Field(i).Tag.Get("db") != "" {
			key = t.Field(i).Tag.Get("db")
		}
		keys = append(keys, key)
		vals = append(vals, v.Field(i).Interface())
	}
	return keys, vals
}

func useIgnore(s any, ignore []string) ([]string, []any) {
	ks, vs := s2d(s)

	var k = make([]string, 0, len(ks))
	var v = make([]any, 0, len(vs))

	for i := 0; i < len(ks); i++ {
		if !slices.Contains(ignore, ks[i]) {
			k = append(k, ks[i])
			v = append(v, vs[i])
		}
	}
	return k, v
}

func useOnly(s any, only []string) ([]string, []any) {
	ks, vs := s2d(s)

	var k = make([]string, 0, len(ks))
	var v = make([]any, 0, len(vs))

	for i := 0; i < len(ks); i++ {
		if slices.Contains(only, ks[i]) {
			k = append(k, ks[i])
			v = append(v, vs[i])
		}
	}
	return k, v
}

func columnScope(s any, opt ...*Opt) ([]string, []any) {
	if len(opt) > 0 {
		var am = opt[0]
		if len(am.only) > 0 {
			return useOnly(s, am.only)
		} else if len(am.ignore) > 0 {
			return useIgnore(s, am.ignore)
		}
	}
	return s2d(s)
}

func whereScope(sb *strings.Builder, opt ...*Opt) {
	if len(opt) > 0 {
		var am = opt[0]
		if len(am.where) > 0 {
			sb.WriteString(" where ")
			sb.WriteString(am.where[0])
			for _, c := range am.where[1:] {
				sb.WriteString(" and ")
				sb.WriteString(c)
			}
		}
	}
}

// Select 使用结构体构建 SQL Select 语句
func Select(table string, s any, opt ...*Opt) string {
	var sb strings.Builder

	var ks, _ = columnScope(s, opt...)

	sb.WriteString("select ")
	sb.WriteString("`")
	sb.WriteString(ks[0])
	sb.WriteString("`")
	for i := 1; i < len(ks); i++ {
		sb.WriteString(",`")
		sb.WriteString(ks[i])
		sb.WriteString("`")
	}
	sb.WriteString(" from `")
	sb.WriteString(table)
	sb.WriteString("`")

	whereScope(&sb, opt...)

	return sb.String()
}

// Update 使用结构体构建 SQL UPDATE 语句
func Update(table string, s any, opt ...*Opt) (string, []any) {
	var data = make([]any, 0)
	var sb strings.Builder

	var ks, vs = columnScope(s, opt...)

	sb.WriteString("update `")
	sb.WriteString(table)
	sb.WriteString("` set ")

	if len(ks) > 0 {
		sb.WriteString("`")
		sb.WriteString(ks[0])
		sb.WriteString("`")
		sb.WriteString(" = ?")
		data = append(data, vs[0])

		for i, k := range ks[1:] {
			sb.WriteString(", `")
			sb.WriteString(k)
			sb.WriteString("`")
			sb.WriteString(" = ?")
			data = append(data, vs[i+1])
		}
	} else {
		panic("no invalid struct")
	}

	whereScope(&sb, opt...)

	return sb.String(), data
}

// Insert 使用结构体构建 SQL INSERT 语句
func Insert(table string, s any, opt ...*Opt) (string, []any) {
	var sb strings.Builder

	var ks, vs = columnScope(s, opt...)

	sb.WriteString("insert into `")
	sb.WriteString(table)
	sb.WriteString("`(")

	sb.WriteString("`")
	sb.WriteString(ks[0])
	sb.WriteString("`")
	for i := 1; i < len(ks); i++ {
		sb.WriteString(",")
		sb.WriteString("`")
		sb.WriteString(ks[i])
		sb.WriteString("`")
	}
	sb.WriteString(") values(")
	sb.WriteString("?")
	for i := 1; i < len(ks); i++ {
		sb.WriteString(",?")
	}
	sb.WriteString(")")
	return sb.String(), vs
}
