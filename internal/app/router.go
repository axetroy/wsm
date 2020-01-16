// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
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
			res.Status = schema.StatusSuccess
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
			hostRouter.GET("", controller.Router(host.GetMyOperationalHostForUser))                                               // 获取我可以操作的服务器信息列表
			hostRouter.POST("", controller.Router(host.CreateHostByUser))                                                         // 创建服务器
			hostRouter.GET("/connection", controller.Router(host.GetHostConnectionRecordListByUser))                              // 获取所有服务器的连接记录
			hostRouter.GET("/connection/:record_id", controller.Router(host.GetHostConnectionRecordDetailByUser))                 // 获取服务器连接记录详情
			hostRouter.PUT("/_/:host_id", controller.Router(host.UpdateHostForUser))                                              // 更新服务器
			hostRouter.GET("/_/:host_id", controller.Router(host.GetHostDetailByIdForUser))                                       // 获取服务器信息
			hostRouter.GET("/_/:host_id/connection", controller.Router(host.GetHostConnectionRecordListByUser))                   // 获取某个服务器的连接记录
			hostRouter.DELETE("/_/:host_id", controller.Router(host.DeleteHostByIdForUser))                                       // 删除服务器
			hostRouter.PUT("/_/:host_id/transfer/:user_id", controller.Router(host.TransferHost))                                 // 转让服务器
			hostRouter.POST("/_/:host_id/collaborator/_/:collaborator_uid", controller.Router(host.AddCollaboratorToHost))        // 添加协作者
			hostRouter.DELETE("/_/:host_id/collaborator/_/:collaborator_uid", controller.Router(host.RemoveCollaboratorFromHost)) // 删除协作者
		}

		// shell 类
		{
			shellRouter := v1.Group("/shell")
			shellRouter.Use(userAuthMiddleware)
			shellRouter.GET("/connect/:host_id", shell.StartTerminalRouter)             // 开启终端，连接 websocket
			shellRouter.GET("/test/:host_id", controller.Router(shell.TestHostConnect)) // 测试服务器是否可连接
			shellRouter.POST("/test", controller.Router(shell.TestPublicServer))        // 测试服务器是否可连接，给定服务器的相关信息即可，无需登陆验证
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
			teamRouter.GET("", controller.Router(team.GetTeamList))                                               // 获取我所在的团队列表
			teamRouter.GET("/all", controller.Router(team.GetAllTeamList))                                        // 获取所有的团队，没有翻页
			teamRouter.POST("", controller.Router(team.CreateTeam))                                               // 创建团队
			teamRouter.GET("/invite", controller.Router(team.GetMyInvitedList))                                   // 获取我的受邀列表
			teamRouter.GET("/_/:team_id", controller.Router(team.GetTeamDetail))                                  // 获取我的团队信息, 只有加入团队才能调用
			teamRouter.PUT("/_/:team_id", controller.Router(team.UpdateTeam))                                     // 更新团队, 只有管理员或者拥有者才能更新
			teamRouter.GET("/_/:team_id/stat", controller.Router(team.StatTeam))                                  // 获取团队的统计信息
			teamRouter.GET("/_/:team_id/connection", controller.Router(host.GetHostConnectionRecordListByTeam))   // 获取团队的服务器连接记录
			teamRouter.GET("/_/:team_id/profile", controller.Router(team.GetMyProfileOfTeam))                     // 获取我在团队中的信息
			teamRouter.GET("/_/:team_id/member/invite", controller.Router(team.GetTeamInviteList))                // 获取发出去的团队邀请
			teamRouter.POST("/_/:team_id/member/invite", controller.Router(team.InviteTeam))                      // 邀请成员加入团队
			teamRouter.PUT("/_/:team_id/member/invite/_/:invite_id", controller.Router(team.ResolveInviteTeam))   // 收邀者 接受/拒绝加入团队
			teamRouter.DELETE("/_/:team_id/member/invite/_/:invite_id", controller.Router(team.CancelInviteTeam)) // 团队管理者取消邀请
			teamRouter.DELETE("/_/:team_id/member/_/:user_id", controller.Router(team.KickUserOutOfTeam))         // 管理员/拥有者 把成员踢出团队
			teamRouter.GET("/_/:team_id/member", controller.Router(team.GetTeamMembers))                          // 获取团队成员列表
			teamRouter.PUT("/_/:team_id/transfer/:user_id", controller.Router(team.TransferTeam))                 // 转让团队
			teamRouter.DELETE("/_/:team_id", controller.Router(team.DeleteTeamByID))                              // 解散团队, 只有拥有者才有权限删除
			teamRouter.DELETE("/_/:team_id/quit", controller.Router(team.QuitTeam))                               // 团队成员退出团队(团队的拥有者无法退出)
			teamRouter.PUT("/_/:team_id/role/:user_id", controller.Router(team.UpdateTeamMemberRole))             // 更改团队成员的角色，只有拥有者和管理员可以操作
			teamRouter.POST("/_/:team_id/host", controller.Router(host.CreateHostByTeam))                         // 添加团队的服务器，只有拥有者和管理员可以操作
			teamRouter.GET("/_/:team_id/host", controller.Router(host.GetMyOperationalHostForTeam))               // 获取团队的服务器列表，只有拥有者和管理员可以操作
			teamRouter.GET("/_/:team_id/host/_/:host_id", controller.Router(host.GetHostDetailByIdForTeam))       // 获取团队的服务器信息，只有拥有者和管理员可以操作
			teamRouter.DELETE("/_/:team_id/host/_/:host_id", controller.Router(host.DeleteHostByIdForTeam))       // 删除服务器，只有拥有者和管理员可以操作
			teamRouter.PUT("/_/:team_id/host/_/:host_id", controller.Router(host.UpdateHostForTeam))              // 更新服务器，只有拥有者和管理员可以操作
		}
	}

	UserRouter = router
}
