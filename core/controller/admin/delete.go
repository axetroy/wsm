// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/helper"
	"github.com/axetroy/terminal/core/middleware"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func DeleteAdminByAccount(account string) {
	database.DeleteRowByTable("admin", "username", account)
}

func DeleteAdminById(c controller.Context, adminId string) (res schema.Response) {
	var (
		err  error
		data schema.AdminProfile
		tx   *gorm.DB
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

		helper.Response(&res, data, err)
	}()

	tx = database.Db.Begin()

	targetAdminInfo := model.Admin{
		Id: adminId,
	}

	if err = tx.First(&targetAdminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
		return
	}

	myInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&myInfo).Error; err != nil {
		return
	}

	// 超级管理员才能操作
	if myInfo.IsSuper == false {
		err = exception.NoPermission
		return
	}

	if err = tx.Delete(model.Admin{
		Id:      targetAdminInfo.Id,
		IsSuper: false, // 超级管理员无法被删除
	}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(targetAdminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	if len(data.Accession) == 0 {
		data.Accession = []string{}
	}

	data.CreatedAt = targetAdminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = targetAdminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func DeleteAdminByIdRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	id := c.Param("admin_id")

	res = DeleteAdminById(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, id)
}
