// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package db

import (
	"log"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
)

var (
	teamID *snowflake.Node
)

func init() {
	node, err := snowflake.NewNode(config.Common.MachineId)

	if err != nil {
		log.Panicln(err)
	}

	teamID = node
}

type Team struct {
	Id      string  `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // 团队ID
	OwnerID string  `gorm:"not null;index;type:varchar(32);" json:"owner_id"`              // 拥有者
	Owner   User    `gorm:"foreignkey:OwnerID" json:"owner"`                               // ** 关联外键 **
	Name    string  `gorm:"not null;type:varchar(32);" json:"name"`                        // 团队名
	Remark  *string `gorm:"null;type:varchar(36);" json:"remark"`                          // 备注

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (t *Team) TableName() string {
	return "team"
}

func (t *Team) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := teamID.Generate().String()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}
	return nil
}
