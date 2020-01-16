// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package user

import (
	"errors"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type SearchUserParams struct {
	Account string `json:"account" form:"account" valid:"required~请输入搜索字段"` // 按用户名来搜索
}

func SearchUser(c *controller.Context) (res schema.Response) {
	var (
		err   error
		input SearchUserParams
		data  = make([]schema.ProfilePublic, 0) // 输出到外部的结果
		list  = make([]db.User, 0)              // 数据库查询出来的原始结果
		tx    *gorm.DB
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

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	userInfo := db.User{}

	if err = db.Db.Model(&userInfo).Limit(100).Where("username LIKE ?", "%"+input.Account+"%").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.ProfilePublic{}
		if err = mapstructure.Decode(v, &d); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}
		data = append(data, d)
	}

	return
}
