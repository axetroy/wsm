// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/middleware"
	"github.com/axetroy/terminal/internal/app/model"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/database"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type Query struct {
	schema.Query
}

func GetList(c controller.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]schema.Profile, 0)
		meta = &schema.Meta{}
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		helper.ResponseList(&res, data, meta, err)
	}()

	query := input.Query

	query.Normalize()

	list := make([]model.User, 0)

	var total int64

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Find(&list).Count(&total).Error; err != nil {
		if err.Error() == exception.EmptyList.Error() {
			err = nil
		} else {
			return
		}
	}

	for _, v := range list {
		d := schema.Profile{}
		if er := mapstructure.Decode(v, &d.ProfilePure); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(list)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

func GetListRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.List{}
		input Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = GetList(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}
