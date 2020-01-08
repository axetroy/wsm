// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package db

import (
	"log"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

var (
	hostConnectionID *snowflake.Node
)

func init() {
	if node, err := snowflake.NewNode(config.Common.MachineId); err != nil {
		log.Panicln(err)
	} else {
		hostConnectionID = node
	}
}

type HostConnectionRecord struct {
	Id      string         `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // 随机生成的主键
	HostID  string         `gorm:"not null;index;type:varchar(32);" json:"host_id"`               // 服务器 ID
	Host    Host           `gorm:"foreignkey:HostID" json:"host"`                                 // **关联外键**
	UserID  string         `gorm:"not null;index;type:varchar(32);" json:"user_id"`               // 操作者的 ID
	User    User           `gorm:"foreignkey:UserID" json:"user"`                                 // **关联外键**
	Records pq.StringArray `gorm:"not null;type:text[]" json:"records"`                           // 操作记录
	Ip      string         `gorm:"not null;type:varchar(32);" json:"ip"`                          // 连接的服务器

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *HostConnectionRecord) TableName() string {
	return "host_connection_record"
}

func (u *HostConnectionRecord) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := hostConnectionID.Generate().String()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}
	return nil
}
