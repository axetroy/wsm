// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package user

import (
	"errors"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/axetroy/wsm/internal/library/util"
	"github.com/axetroy/wsm/internal/library/validator"
	"github.com/jinzhu/gorm"
)

type UpdatePasswordParams struct {
	OldPassword string `json:"old_password" valid:"required~请输入旧密码"`
	NewPassword string `json:"new_password" valid:"required~请输入新密码"`
}

func UpdatePassword(c *controller.Context) (res schema.Response) {
	var (
		err   error
		input UpdatePasswordParams
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

		helper.Response(&res, nil, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	tx = db.Db.Begin()

	userInfo := db.User{Id: c.Uid}

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

	if err = tx.Model(&userInfo).Update(db.User{Password: newPassword}).Error; err != nil {
		return
	}

	return
}
