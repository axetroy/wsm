// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/helper"
	"github.com/axetroy/terminal/core/middleware"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/rbac/accession"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/axetroy/terminal/core/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateParams struct {
	Status    *model.AdminStatus `json:"status"`    // 管理员状态
	Name      *string            `json:"name"`      // 管理员名字
	Accession *[]string          `json:"accession"` // 管理员的权限
}

func Update(c controller.Context, adminId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.AdminProfile
		tx           *gorm.DB
		shouldUpdate bool
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
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	myInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&myInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	if !myInfo.IsSuper {
		err = exception.AdminNotSuper
		return
	}

	adminInfo := model.Admin{
		Id: adminId,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
		return
	}

	// 需要更新的字段
	updated := map[string]interface{}{}

	if input.Status != nil {
		shouldUpdate = true
		adminInfo.Status = *input.Status
		updated["status"] = *input.Status
	}

	if input.Name != nil {
		shouldUpdate = true
		adminInfo.Name = *input.Name
		updated["name"] = *input.Name
	}

	if input.Accession != nil {
		shouldUpdate = true
		accessions := accession.FilterAdminAccession(*input.Accession) // 提取有效的权限, 无效的忽略
		adminInfo.Accession = accessions
		updated["accession"] = accessions
	}

	if shouldUpdate {
		if err = tx.Model(&adminInfo).Updates(updated).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.AdminNotExist
				return
			}
			return
		}
	}

	if err = mapstructure.Decode(adminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	if len(data.Accession) == 0 {
		data.Accession = []string{}
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	id := c.Param("admin_id")

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, id, input)
}
