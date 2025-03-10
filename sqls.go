package sqls

import (
	"reflect"
	"strings"
)

// 结构体转换为数据
// 这里不转换为map而转换为切片，避免每次顺序不一样
func s2d(s any, ignore ...string) ([]string, []any) {
	var keys []string
	var vals []any

	var t = reflect.TypeOf(s)
	var v = reflect.ValueOf(s)
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
		for _, ign := range ignore { // 跳过忽略的字段
			if key == ign {
				goto next
			}
		}
		keys = append(keys, key)
		vals = append(vals, v.Field(i).Interface())
	next:
	}
	return keys, vals
}

// Select 使用结构体构建 SQL Select 语句
func Select(table string, s any, applies ...Apply) string {
	var sb strings.Builder

	var am = newApplyMap(applies...)
	ks, _ := s2d(s, am.ignore()...)

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

	am.where(&sb)

	return sb.String()
}

// Update 使用结构体构建 SQL UPDATE 语句
func Update(table string, s any, applies ...Apply) (string, []any) {
	var data = make([]any, 0)
	var sb strings.Builder

	var am = newApplyMap(applies...)
	ks, vs := s2d(s, am.ignore()...)

	sb.WriteString("update `")
	sb.WriteString(table)
	sb.WriteString("` set ")

	if len(ks) > 0 {
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

	am.where(&sb)

	return sb.String(), data
}

// Insert 使用结构体构建 SQL INSERT 语句
func Insert(table string, s any, applies ...Apply) (string, []any) {
	var sb strings.Builder

	var am = newApplyMap(applies...)
	ks, vs := s2d(s, am.ignore()...)

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
