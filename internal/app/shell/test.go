// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package shell

import (
	"errors"

	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/axetroy/wsm/internal/library/session"
)

type TestPublicServerParams struct {
	Username string `json:"username" valid:"required~请输入用户名"`
	Host     string `json:"host" valid:"required~请输入服务器,host~请输入正确的服务器地址"`
	Port     uint   `json:"port" valid:"required~请输入端口,port~请输入正确的端口,range(1|65535)"`
	Password string `json:"password" valid:"required~请输入密码"`
}

func TestPublicServer(c *controller.Context) (res schema.Response) {
	var (
		err   error
		input TestPublicServerParams
		data  bool
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

		helper.Response(&res, data, nil, err)
	}()

	if err = c.Validator(input); err != nil {
		return
	}

	config := session.Config{
		Username: input.Username,
		Host:     input.Host,
		Port:     input.Port,
		Password: input.Password,
	}

	data = session.Test(config)

	return
}
