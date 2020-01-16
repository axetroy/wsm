// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package app

import (
	"fmt"
	"net/http"
	"path"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/host"
	"github.com/axetroy/wsm/internal/app/middleware"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/app/shell"
	"github.com/axetroy/wsm/internal/app/team"
	"github.com/axetroy/wsm/internal/app/user"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/dotenv"
	"github.com/gin-gonic/gin"
)

var UserRouter *gin.Engine

func init() {
	if config.Common.Mode == config.ModeProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

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

		v1.GET("", controller.Router(func(c *controller.Context) (res schema.Response) {
			res.Data = "pong"
			return
		}))

		userAuthMiddleware := middleware.Authenticate(false) // 用户 Token 的中间件

		// 认证类
		{
			authRouter := v1.Group("/auth")
			authRouter.POST("/signup", controller.Router(user.SignUpWithUsername)) // 注册账号, 通过用户名+密码
			authRouter.POST("/signin", controller.Router(user.LoginWithUsername))  // 登陆账号
		}

		// 服务器管理
		{
			hostRouter := v1.Group("/host")
			hostRouter.Use(userAuthMiddleware)
			hostRouter.GET("", host.Core.QueryMyOperationalServerRouter)                                                  // 获取我可以操作的服务器信息列表
			hostRouter.POST("", host.Core.CreateHostByUserRouter)                                                         // 创建服务器
			hostRouter.GET("/connection", host.Core.QueryHostConnectionRecordListRouter)                                  // 获取所有服务器的连接记录
			hostRouter.GET("/connection/:record_id", host.Core.QueryHostConnectionRecordRouter)                           // 获取服务器连接记录详情
			hostRouter.PUT("/_/:host_id", host.Core.UpdateHostRouter)                                                     // 更新服务器
			hostRouter.GET("/_/:host_id", host.Core.QueryMyHostByIDRouter)                                                // 获取服务器信息
			hostRouter.GET("/_/:host_id/connection", host.Core.QueryHostConnectionRecordListRouter)                       // 获取某个服务器的连接记录
			hostRouter.DELETE("/_/:host_id", host.Core.DeleteHostByIDRouter)                                              // 删除服务器
			hostRouter.PUT("/_/:host_id/transfer/:user_id", host.Core.TransferHostRouter)                                 // 转让服务器
			hostRouter.POST("/_/:host_id/collaborator/_/:collaborator_uid", host.Core.AddCollaboratorToHostRouter)        // 添加协作者
			hostRouter.DELETE("/_/:host_id/collaborator/_/:collaborator_uid", host.Core.RemoveCollaboratorFromHostRouter) // 删除协作者
		}

		// shell 类
		{
			shellRouter := v1.Group("/shell")
			shellRouter.Use(userAuthMiddleware)
			shellRouter.GET("/connect/:host_id", shell.Core.StartTerminalRouter) // 开启终端，连接 websocket
			shellRouter.GET("/test/:host_id", shell.Core.TestHostConnectRouter)  // 测试服务器是否可连接
			shellRouter.POST("/test", shell.Core.TestPublicServerRouter)         // 测试服务器是否可连接，给定服务器的相关信息即可，无需登陆验证
		}

		// 用户类
		{
			userRouter := v1.Group("/user")
			userRouter.Use(userAuthMiddleware)
			userRouter.GET("/profile", controller.Router(user.GetProfile))      // 获取用户详细信息
			userRouter.PUT("/profile", controller.Router(user.UpdateProfile))   // 更新用户资料
			userRouter.PUT("/password", controller.Router(user.UpdatePassword)) // 更新登陆密码
			userRouter.GET("/search", controller.Router(user.SearchUser))       // 搜索用户
		}

		// 团队相关
		{
			teamRouter := v1.Group("/team")
			teamRouter.Use(userAuthMiddleware)
			teamRouter.GET("", team.Core.QueryMyTeamsRouter)                                              // 获取我所在的团队列表
			teamRouter.GET("/all", team.Core.GetAllTeamsRouter)                                           // 获取所有的团队，没有翻页
			teamRouter.POST("", team.Core.CreateTeamRouter)                                               // 创建团队
			teamRouter.GET("/invite", team.Core.GetMyInvitedRecordRouter)                                 // 获取我的受邀列表
			teamRouter.GET("/_/:team_id", team.Core.QueryMyTeamRouter)                                    // 获取我的团队信息, 只有加入团队才能调用
			teamRouter.PUT("/_/:team_id", team.Core.UpdateTeamRouter)                                     // 更新团队, 只有管理员或者拥有者才能更新
			teamRouter.GET("/_/:team_id/stat", team.Core.StatTeamRouter)                                  // 获取团队的统计信息
			teamRouter.GET("/_/:team_id/connection", host.Core.QueryTeamHostConnectionRecordListRouter)   // 获取团队的服务器连接记录
			teamRouter.GET("/_/:team_id/profile", team.Core.GetMyProfileRouter)                           // 获取我在团队中的信息
			teamRouter.GET("/_/:team_id/member/invite", team.Core.GetTeamInviteRecordRouter)              // 获取发出去的团队邀请
			teamRouter.POST("/_/:team_id/member/invite", team.Core.InviteTeamRouter)                      // 邀请成员加入团队
			teamRouter.PUT("/_/:team_id/member/invite/_/:invite_id", team.Core.ResolveInviteTeamRouter)   // 收邀者 接受/拒绝加入团队
			teamRouter.DELETE("/_/:team_id/member/invite/_/:invite_id", team.Core.CancelInviteTeamRouter) // 团队管理者取消邀请
			teamRouter.DELETE("/_/:team_id/member/_/:user_id", team.Core.KickOutByUIDRouter)              // 管理员/拥有者 把成员踢出团队
			teamRouter.GET("/_/:team_id/member", team.Core.QueryTeamMembersRouter)                        // 获取团队成员列表
			teamRouter.PUT("/_/:team_id/transfer/:user_id", team.Core.TransferTeamRouter)                 // 转让团队
			teamRouter.DELETE("/_/:team_id", team.Core.DeleteTeamByIDRouter)                              // 解散团队, 只有拥有者才有权限删除
			teamRouter.DELETE("/_/:team_id/quit", team.Core.QuitTeamRouter)                               // 团队成员退出团队(团队的拥有者无法退出)
			teamRouter.PUT("/_/:team_id/role/:user_id", team.Core.UpdateMemberRoleRouter)                 // 更改团队成员的角色，只有拥有者和管理员可以操作
			teamRouter.POST("/_/:team_id/host", host.Core.CreateHostByTeamRouter)                         // 添加团队的服务器，只有拥有者和管理员可以操作
			teamRouter.GET("/_/:team_id/host", host.Core.QueryHostByTeamRouter)                           // 获取团队的服务器列表，只有拥有者和管理员可以操作
			teamRouter.GET("/_/:team_id/host/_/:host_id", host.Core.QueryMyHostByTeamRouter)              // 获取团队的服务器信息，只有拥有者和管理员可以操作
			teamRouter.DELETE("/_/:team_id/host/_/:host_id", host.Core.DeleteHostByTeamRouter)            // 删除服务器，只有拥有者和管理员可以操作
			teamRouter.PUT("/_/:team_id/host/_/:host_id", host.Core.UpdateHostByTeamRouter)               // 更新服务器，只有拥有者和管理员可以操作
		}
	}

	UserRouter = router
}
