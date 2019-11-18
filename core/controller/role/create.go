// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/model"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/database"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/axetroy/terminal/internal/library/rbac/accession"
	"github.com/axetroy/terminal/internal/library/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateParams struct {
	Name        string   `json:"name" valid:"required~请输入角色名"`       // 角色名
	Description string   `json:"description" valid:"required~请输入描述"` // 描述
	Accession   []string `json:"accession" valid:"required~请输入权限"`   // 权限列表
	Note        *string  `json:"note"`                               // 备注
}

func Create(c controller.Context, input CreateParams) (res schema.Response) {
	var (
		err  error
		data schema.Role
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

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	if accession.Valid(input.Accession) == false {
		err = exception.InvalidParams
		return
	}

	roleInfo := model.Role{
		Name:        input.Name,
		Description: input.Description,
		Accession:   input.Accession,
	}

	if err = tx.Create(&roleInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(roleInfo, &data.RolePure); er != nil {
		err = er
		return
	}

	data.CreatedAt = roleInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = roleInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateRouter(c *gin.Context) {
	var (
		input CreateParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Create(controller.NewContext(c), input)
}
