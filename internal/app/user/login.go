// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package user

import (
	"errors"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/axetroy/wsm/internal/library/password"
	"github.com/axetroy/wsm/internal/library/token"
	"github.com/axetroy/wsm/internal/library/validator"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type SignInParams struct {
	Account  string `json:"account" valid:"required~请输入登陆账号"`
	Password string `json:"password" valid:"required~请输入密码"`
}

func LoginWithUsername(c *controller.Context) (res schema.Response) {
	var (
		err   error
		input SignInParams
		data  = &schema.ProfileWithToken{}
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

	if validator.IsPhone(input.Account) {
		// 用手机号登陆
		userInfo.Phone = &input.Account
	} else if validator.IsEmail(input.Account) {
		// 用邮箱登陆
		userInfo.Email = &input.Account
	} else {
		// 用用户名
		userInfo.Username = input.Account
	}

	tx = db.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	if password.Verify(input.Password, userInfo.Password) == false {
		err = exception.InvalidAccountOrPassword
		return
	}

	if err = userInfo.CheckStatusValid(); err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(config.Http.Secret, userInfo.Id); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	return
}
