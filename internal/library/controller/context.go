// Copyright 2019 Axetroy. All rights reserved. MIT license.
package controller

import (
	"github.com/axetroy/terminal/internal/app/middleware"
	"github.com/axetroy/terminal/internal/library/validator"
	"github.com/gin-gonic/gin"
)

// 控制器的上下文
type Context struct {
	*gin.Context
	Uid       string `json:"uid"`        // 操作人的用户 ID
	UserAgent string `json:"user_agent"` // 用户代理
	Ip        string `json:"ip"`         // IP地址
}

func (c *Context) Validator(input interface{}) error {
	if err := validator.ValidateStruct(input); err != nil {
		return err
	}
	return nil
}

func NewContextFromGinContext(c *gin.Context) Context {
	return Context{
		Uid:       c.GetString(middleware.ContextUidField),
		UserAgent: c.GetHeader("user-agent"),
		Ip:        c.ClientIP(),
	}
}

func NewContext(uid string, userAgent string, ip string) Context {
	return Context{
		Uid:       uid,
		UserAgent: userAgent,
		Ip:        ip,
	}
}
