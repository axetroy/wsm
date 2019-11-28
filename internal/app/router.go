// Copyright 2019 Axetroy. All rights reserved. MIT license.
package app

import (
	"fmt"
	"net/http"
	"path"

	"github.com/axetroy/terminal/internal/app/config"
	"github.com/axetroy/terminal/internal/app/host"
	"github.com/axetroy/terminal/internal/app/middleware"
	"github.com/axetroy/terminal/internal/app/oauth2"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/app/shell"
	"github.com/axetroy/terminal/internal/app/user"
	"github.com/axetroy/terminal/internal/library/dotenv"
	"github.com/gin-gonic/gin"
)

var UserRouter *gin.Engine

func init() {
	if config.Common.Mode == config.ModeProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.LoadHTMLGlob("view/*")
	router.StaticFS("/static", http.Dir("static"))

	router.Use(middleware.GracefulExit())

	router.Use(middleware.CORS())

	router.Static("/public", path.Join(dotenv.RootDir, "public"))

	if config.Common.Mode != config.ModeProduction {
		router.Use(gin.Logger())
	}

	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, schema.Response{
			Status:  schema.StatusFail,
			Message: fmt.Sprintf("%v ", http.StatusNotFound) + http.StatusText(http.StatusNotFound),
			Data:    nil,
		})
	})

	{
		v1 := router.Group("/v1")
		v1.Use(middleware.Common)

		v1.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ping": "pong"})
		})

		userAuthMiddleware := middleware.Authenticate(false) // 用户Token的中间件

		// 认证类
		{
			authRouter := v1.Group("/auth")
			authRouter.POST("/signup", user.Core.SignUpWithUsernameRouter) // 注册账号, 通过用户名+密码
			authRouter.POST("/signin", user.Core.LoginWithUsernameRouter)  // 登陆账号
		}

		// host 类
		{
			hostRouter := v1.Group("/host")
			hostRouter.GET("", userAuthMiddleware, host.Core.QueryOperationalServerRouter)                                                    // 获取我可以操作的服务器信息列表
			hostRouter.POST("", userAuthMiddleware, host.Core.CreateHostRouter)                                                               // 创建服务器
			hostRouter.PUT("/_/:host_id", userAuthMiddleware, host.Core.UpdateHostRouter)                                                     // 更新服务器
			hostRouter.GET("/_/:host_id", userAuthMiddleware, host.Core.QueryHostRouter)                                                      // 获取服务器信息
			hostRouter.DELETE("/_/:host_id", userAuthMiddleware, host.Core.QueryOperationalServerRouter)                                      // 删除服务器
			hostRouter.PUT("/_/:host_id/transfer/:user_id", userAuthMiddleware, host.Core.TransferHostRouter)                                 // 转让服务器
			hostRouter.POST("/_/:host_id/collaborator/_/:collaborator_uid", userAuthMiddleware, host.Core.AddCollaboratorToHostRouter)        // 添加协作者
			hostRouter.DELETE("/_/:host_id/collaborator/_/:collaborator_uid", userAuthMiddleware, host.Core.RemoveCollaboratorFromHostRouter) // 删除协作者
		}

		// shell 类
		{
			shellRouter := v1.Group("/shell")
			shellRouter.GET("/demo", shell.Core.ExampleRouter)
			shellRouter.GET("", userAuthMiddleware, shell.Core.StartTerminalRouter)
		}

		// oAuth2 认证
		{
			oAuthRouter := v1.Group("/oauth2")
			oAuthRouter.GET("/:provider", oauth2.Core.AuthRouter)              // 前去进行 oAuth 认证
			oAuthRouter.GET("/:provider/callback", oauth2.Core.CallbackRouter) // 认证成功后，跳转回来的回调地址
		}

		// 用户类
		{
			userRouter := v1.Group("/user")
			userRouter.Use(userAuthMiddleware)                          // 用户登出
			userRouter.GET("/profile", user.Core.GetProfileRouter)      // 获取用户详细信息
			userRouter.PUT("/profile", user.Core.UpdateProfileRouter)   // 更新用户资料
			userRouter.PUT("/password", user.Core.UpdatePasswordRouter) // 更新登陆密码
		}
	}

	UserRouter = router
}
