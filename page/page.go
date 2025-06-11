package page

import "fmt"

type Pager struct {
	Size int64 `form:"size,default=10"` // 分页大小
	Page int64 `form:"page,default=1"`  // 页码
}

func (p *Pager) Limit() int64 {
	if p.Size < 1 || p.Size > 1000 {
		p.Size = 10 // 校正
	}
	return p.Size
}

func (p *Pager) Offset() int64 {
	if p.Page > 0 {
		return (p.Page - 1) * p.Limit()
	}
	return 0
}

func (p *Pager) Pagination(total int64) *Pagination {
	return &Pagination{
		Total: total,
		Page:  p.Page,
		Size:  p.Size,
	}
}

type Sorter struct {
	Field string `form:"field,optional"` // 字段
	Asc   bool   `form:"asc,optional"`   // 是否正序
}

func (s *Sorter) OrderBy(f ...func(string) string) string {
	if s.Field != "" {
		var field string
		if len(f) > 0 {
			field = f[0](s.Field)
			if field == "" { // 被f过滤掉了
				return ""
			}
		} else {
			field = s.Field
		}
		if s.Asc {
			return fmt.Sprintf("order by %s", field)
		} else {
			return fmt.Sprintf("order by %s desc", field)
		}
	}
	return ""
}

type Pagination struct {
	Size  int64 `json:"size"`  // 分页数量
	Page  int64 `json:"page"`  // 页码
	Total int64 `json:"total"` // 总条数
}

type PaginationRsp[T any] struct {
	Items      []T         `json:"items"`      // 数据列表
	Pagination *Pagination `json:"pagination"` // 分页数据
}
