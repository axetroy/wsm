// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package schema

import (
	"fmt"
	"regexp"
	"strings"
)

type Order string

type Query struct {
	Limit    int     `json:"limit" form:"limit"`
	Page     int     `json:"page" form:"page"`
	Sort     string  `json:"sort" form:"sort"`
	Platform *string `json:"platform" form:"platform"`
}

type Sort struct {
	Field string `json:"field"` // 排序的字段
	Order Order  `json:"order"` // 字段的排序方向
}

var (
	DefaultLimit       = 10            // 默认只获取 10 条数据
	DefaultPage        = 1             // 默认第 1 页
	DefaultSort        = "-created_at" // 默认按照创建时间排序
	MaxLimit           = 100           // 最大的查询数量，100 条 防止查询数据过大拖慢服务端性能
	OrderAsc     Order = "ASC"         // 排序方式，正序
	OrderDesc    Order = "DESC"        // 排序方式，倒序
	ascReg             = regexp.MustCompile("^\\+")
	descReg            = regexp.MustCompile("^-")
)

func (q *Query) formatSort() (fields []Sort) {
	arr := strings.Split(q.Sort, ",")

	for _, field := range arr {
		s := strings.Split(field, "")

		switch s[0] {

		case "-":
			fields = append(fields, Sort{
				Field: descReg.ReplaceAllString(field, ""),
				Order: OrderDesc,
			})
		default:
			fields = append(fields, Sort{
				Field: ascReg.ReplaceAllString(field, ""),
				Order: OrderAsc,
			})
		}
	}

	return
}

func (q *Query) orderStr() []string {
	sorts := q.formatSort()

	arr := make([]string, 0)

	for _, field := range sorts {
		arr = append(arr, fmt.Sprintf("%s %s", field.Field, field.Order))
	}

	return arr
}

func (q *Query) Order() string {
	return strings.Join(q.orderStr(), ", ")
}

func (q *Query) Offset() uint32 {
	return uint32((q.Page - 1) * q.Limit)
}

func (q *Query) Normalize() *Query {

	if q.Limit <= 0 {
		q.Limit = DefaultLimit // 默认查询10条
	} else if q.Limit > MaxLimit {
		q.Limit = MaxLimit
	}

	if q.Page <= 0 {
		q.Page = DefaultPage
	}

	if q.Sort == "" {
		q.Sort = DefaultSort
	}

	return q
}
