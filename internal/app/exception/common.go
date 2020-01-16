// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package exception

var (
	SystemMaintenance = New("系统维护中", -1)
	Unknown           = New("未知错误", 0)
	DataBase          = New("数据库错误", 100)
	DataBinding       = New("数据转换错误", 101)
	InvalidParams     = New("参数不正确", 100000)
	NoData            = New("找不到数据", 100001)
	NoPermission      = New("没有权限", 100002)
	InvalidFormat     = New("格式不正确", 100003)
	Duplicate         = New("重复", 100004)
	InvalidAuth       = New("无效的身份认证方式", 999999)
	InvalidToken      = New("无效的身份令牌", 999999)
	TokenExpired      = New("身份令牌已过期", 999999)

	// 用户类
	UserNotExist             = New("用户不存在", 200000)
	UserExist                = New("用户已存在", 200001)
	UserIsInActive           = New("帐号未激活", 200003)
	UserHaveBeenBan          = New("帐号已被禁用", 200004)
	PasswordDuplicate        = New("新密码和旧密码不能相同", 200005)
	InvalidAccountOrPassword = New("账号或密码错误", 200006)
	InvalidPassword          = New("密码错误", 200012)
)
