// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"net/http"

	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/helper"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/axetroy/terminal/core/util"
	"github.com/axetroy/terminal/core/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UpdatePasswordParams struct {
	OldPassword string `json:"old_password" valid:"required~请输入旧密码"`
	NewPassword string `json:"new_password" valid:"required~请输入新密码"`
}

func UpdatePassword(c controller.Context, input UpdatePasswordParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

		helper.Response(&res, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	tx = database.Db.Begin()

	userInfo := model.User{Id: c.Uid}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 验证密码是否正确
	if userInfo.Password != util.GeneratePassword(input.OldPassword) {
		err = exception.InvalidPassword
		return
	}

	newPassword := util.GeneratePassword(input.NewPassword)

	if err = tx.Model(&userInfo).Update(model.User{Password: newPassword}).Error; err != nil {
		return
	}

	return
}

func UpdatePasswordRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdatePasswordParams
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

	res = UpdatePassword(controller.NewContext(c), input)
}
