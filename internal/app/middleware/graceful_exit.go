// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package middleware

import (
	schema2 "github.com/axetroy/wsm/internal/app/schema"
	"net/http"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/gin-gonic/gin"
)

// 优雅退出中间件
// 再接收到退出指令之后，则 HTTP 服务不再接收新的请求
func GracefulExit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Common.Exiting {
			err := exception.SystemMaintenance
			c.JSON(http.StatusOK, schema2.Response{
				Status:  err.Code(),
				Message: err.Error(),
				Data:    nil,
			})
			c.Abort()
		}
	}
}
