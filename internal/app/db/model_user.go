// Copyright 2019 Axetroy. All rights reserved. Apache License 2.0.
package db

import (
	"log"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type UserStatus int32

type Gender int

const (
	// 用户状态
	UserStatusBanned      UserStatus = -100 // 账号被禁用
	UserStatusInactivated UserStatus = -1   // 账号未激活
	UserStatusInit        UserStatus = 1    // 初始化状态

	// 用户性别
	GenderUnknown Gender = 0 // 未知性别
	GenderMale               // 男
	GenderFemale             // 女
)

var (
	userID *snowflake.Node
)

func init() {
	node, err := snowflake.NewNode(config.Common.MachineId)

	if err != nil {
		log.Panicln(err)
	}

	userID = node
}

type User struct {
	Id       string         `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 用户ID
	Username string         `gorm:"not null;type:varchar(36)unique;index" json:"username"`        // 用户名
	Password string         `gorm:"not null;type:varchar(36);index" json:"password"`              // 登陆密码
	Nickname *string        `gorm:"null;type:varchar(36)" json:"nickname"`                        // 昵称
	Phone    *string        `gorm:"null;unique;type:varchar(16);index" json:"phone"`              // 手机号
	Email    *string        `gorm:"null;unique;type:varchar(36);index" json:"email"`              // 邮箱
	Status   UserStatus     `gorm:"not null" json:"status"`                                       // 状态
	Role     pq.StringArray `gorm:"not null;type:varchar(36)[]" json:"role"`                      // 角色, 用户可以拥有多个角色
	Avatar   string         `gorm:"not null;type:varchar(128)" json:"avatar"`                     // 头像
	Gender   Gender         `gorm:"default(0)" json:"gender"`                                     // 性别

	// 外键关联
	OAuth []OAuth `gorm:"foreignkey:Uid" json:"oauth"` // **外键**

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := userID.Generate().String()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}
	return nil
}

// 检查用户状态是否正常
func (u *User) CheckStatusValid() error {
	if u.Status != UserStatusInit {
		switch u.Status {
		case UserStatusInactivated:
			return exception.UserIsInActive
		case UserStatusBanned:
			return exception.UserHaveBeenBan
		}
	}

	return nil
}
