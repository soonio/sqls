package sqls

import "strings"

type Apply func() (string, []string)

// ApplyIgnore 配置忽律某些字段
func ApplyIgnore(column ...string) Apply {
	return func() (string, []string) {
		return "ignore", column
	}
}

// ApplyWhere 配置使用where条件
func ApplyWhere(condition ...string) Apply {
	return func() (string, []string) {
		return "where", condition
	}
}

type applyMap map[string][]string

func newApplyMap(applies ...Apply) applyMap {
	var am = make(applyMap)
	for _, apply := range applies {
		k, v := apply()
		am[k] = v
	}
	return am
}

func (m applyMap) where(sb *strings.Builder) {
	if where, ok := m[`where`]; ok {
		if len(where) > 0 {
			sb.WriteString(" where ")
			sb.WriteString(where[0])
			for _, c := range where[1:] {
				sb.WriteString(" and ")
				sb.WriteString(c)
			}
		}
	}
}

func (m applyMap) ignore() []string {
	if ignore, ok := m[`ignore`]; ok {
		return ignore
	}
	return nil
}
