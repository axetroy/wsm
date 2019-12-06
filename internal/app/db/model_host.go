// Copyright 2019 Axetroy. All rights reserved. MIT license.
package db

import (
	"time"

	"github.com/axetroy/terminal/internal/library/util"
	"github.com/jinzhu/gorm"
)

type HostOwnerType string

const (
	HostOwnerTypeUser HostOwnerType = "user"
	HostOwnerTypeTeam HostOwnerType = "team"
)

type Host struct {
	Id         string        `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // 用户ID
	OwnerID    string        `gorm:"not null;index;type:varchar(32);" json:"owner_id"`              // 拥有者ID，可以是userID，也可以是 teamID
	User       User          `gorm:"foreignkey:OwnerID" json:"user"`                                // **关联外键**
	Team       Team          `gorm:"foreignkey:OwnerID" json:"team"`                                // **关联外键**
	OwnerType  HostOwnerType `gorm:"not null;index;type:varchar(16);" json:"owner_type"`            // 拥有者的类型，是个人用户拥有还是组织拥有?
	Name       string        `gorm:"not null;type:varchar(32);" json:"name"`                        // 服务器名
	Host       string        `gorm:"not null;type:varchar(36);index;" json:"host"`                  // 服务器地址
	Port       uint          `gorm:"not null;index;" json:"port"`                                   // 端口
	Username   string        `gorm:"not null;type:varchar(36);index;" json:"username"`              // 用户名
	Password   string        `gorm:"not null;type:varchar(36);index;" json:"password"`              // 登陆密码
	PrivateKey string        `gorm:"not null;type:varchar(36);index;" json:"private_key"`           // 登陆的私钥
	Remark     *string       `gorm:"null;type:varchar(36);" json:"remark"`                          // 备注

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *Host) TableName() string {
	return "host"
}

func (u *Host) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := util.GenerateId()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}
	return nil
}
