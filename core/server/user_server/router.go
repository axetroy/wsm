// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user_server

import (
	"fmt"
	"net/http"
	"path"

	"github.com/axetroy/terminal/core/config"
	"github.com/axetroy/terminal/core/controller/auth"
	"github.com/axetroy/terminal/core/controller/downloader"
	"github.com/axetroy/terminal/core/controller/oauth2"
	"github.com/axetroy/terminal/core/controller/resource"
	"github.com/axetroy/terminal/core/controller/shell"
	"github.com/axetroy/terminal/core/controller/uploader"
	"github.com/axetroy/terminal/core/controller/user"
	"github.com/axetroy/terminal/core/middleware"
	"github.com/axetroy/terminal/core/rbac"
	"github.com/axetroy/terminal/core/rbac/accession"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/dotenv"
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
			authRouter.POST("/signup", auth.SignUpWithUsernameRouter)     // 注册账号, 通过用户名+密码
			authRouter.POST("/signin/oauth2", auth.SignInWithOAuthRouter) // oAuth 码登陆
			authRouter.POST("/signin", auth.SignInRouter)                 // 登陆账号
		}

		{
			shellRouter := v1.Group("/shell")
			shellRouter.GET("/demo", shell.DemoRouter)
			shellRouter.GET("", shell.StartRouter)
		}

		// oAuth2 认证
		{
			oAuthRouter := v1.Group("/oauth2")
			oAuthRouter.GET("/:provider", oauth2.AuthRouter)                  // 前去进行 oAuth 认证
			oAuthRouter.GET("/:provider/callback", oauth2.AuthCallbackRouter) // 认证成功后，跳转回来的回调地址
		}

		// 用户类
		{
			userRouter := v1.Group("/user")
			userRouter.Use(userAuthMiddleware)
			userRouter.GET("/logout", user.SignOut)                                                         // 用户登出
			userRouter.GET("/profile", user.GetProfileRouter)                                               // 获取用户详细信息
			userRouter.PUT("/profile", rbac.Require(*accession.ProfileUpdate), user.UpdateProfileRouter)    // 更新用户资料
			userRouter.PUT("/password", rbac.Require(*accession.PasswordUpdate), user.UpdatePasswordRouter) // 更新登陆密码
		}

		// 通用类
		{
			// 文件上传
			v1.POST("/upload/file", uploader.File)      // 上传文件
			v1.POST("/upload/image", uploader.Image)    // 上传图片
			v1.GET("/upload/example", uploader.Example) // 上传文件的 example
			// 单纯获取资源文本
			v1.GET("/resource/file/:filename", resource.File)           // 获取文件纯文本
			v1.GET("/resource/image/:filename", resource.Image)         // 获取图片纯文本
			v1.GET("/resource/thumbnail/:filename", resource.Thumbnail) // 获取缩略图纯文本
			// 下载资源
			v1.GET("/download/file/:filename", downloader.File)           // 下载文件
			v1.GET("/download/image/:filename", downloader.Image)         // 下载图片
			v1.GET("/download/thumbnail/:filename", downloader.Thumbnail) // 下载缩略图
			// 公共资源目录
			v1.GET("/avatar/:filename", user.GetAvatarRouter) // 获取用户头像
		}

	}

	UserRouter = router
}
