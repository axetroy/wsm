// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package db

import (
	"log"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
)

type HostOwnerType string
type HostConnectType string

const (
	HostOwnerTypeUser         HostOwnerType   = "user"
	HostOwnerTypeTeam         HostOwnerType   = "team"
	HostConnectTypePassword   HostConnectType = "password"
	HostConnectTypePrivateKey HostConnectType = "private_key"
)

var (
	hostID *snowflake.Node
)

func init() {
	node, err := snowflake.NewNode(config.Common.MachineId)

	if err != nil {
		log.Panicln(err)
	}

	hostID = node
}

type Host struct {
	Id          string          `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // 用户ID
	OwnerID     string          `gorm:"not null;index;type:varchar(32);" json:"owner_id"`              // 拥有者ID，可以是userID，也可以是 teamID
	User        User            `gorm:"foreignkey:OwnerID" json:"user"`                                // **关联外键**
	Team        Team            `gorm:"foreignkey:OwnerID" json:"team"`                                // **关联外键**
	OwnerType   HostOwnerType   `gorm:"not null;index;type:varchar(16);" json:"owner_type"`            // 拥有者的类型，是个人用户拥有还是组织拥有?
	Name        string          `gorm:"not null;type:varchar(32);" json:"name"`                        // 服务器名
	Host        string          `gorm:"not null;type:varchar(36);index;" json:"host"`                  // 服务器地址
	Port        uint            `gorm:"not null;index;" json:"port"`                                   // 端口
	Username    string          `gorm:"not null;type:varchar(36);index;" json:"username"`              // 用户名
	ConnectType HostConnectType `gorm:"not null;type:varchar(36);index;" json:"connect_type"`          // 服务器的连接方式
	Passport    string          `gorm:"not null;type:varchar(4096);" json:"passport"`                  // 密码/密钥
	Remark      *string         `gorm:"null;type:varchar(36);" json:"remark"`                          // 备注

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *Host) TableName() string {
	return "host"
}

func (u *Host) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := hostID.Generate().String()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}
	return nil
}
