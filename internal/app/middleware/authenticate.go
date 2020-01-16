// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package middleware

import (
	"net/http"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/token"
	"github.com/gin-gonic/gin"
)

var (
	ContextUidField = "uid"
)

// Token 验证中间件
func Authenticate(isAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err         error
			tokenString string
			status      = schema.StatusFail
		)
		defer func() {
			if err != nil {
				c.JSON(http.StatusOK, schema.Response{
					Status:  status,
					Message: err.Error(),
					Data:    nil,
				})
				c.Abort()
			}
		}()

		if s, isExist := c.GetQuery(token.AuthField); isExist == true {
			tokenString = s
		} else {
			tokenString = c.GetHeader(token.AuthField)

			if len(tokenString) == 0 {
				if s, er := c.Cookie(token.AuthField); er != nil {
					err = exception.InvalidToken
					status = exception.InvalidToken.Code()
					return
				} else {
					tokenString = s
				}
			}
		}

		if claims, er := token.Parse(config.Http.Secret, tokenString); er != nil {
			err = er
			status = exception.InvalidToken.Code()
			return
		} else {
			// 把 UID 挂载到上下文中国呢
			c.Set(ContextUidField, claims.Uid)
		}
	}
}
