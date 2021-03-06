// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package db

import (
	"fmt"
	"log"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	Db     *gorm.DB
	Config = config.Database
)

func init() {
	DataSourceName := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", Config.Driver, Config.Username, Config.Password, Config.Host, Config.Port, Config.DatabaseName)

	log.Println("正在连接数据库...")

	db, err := gorm.Open(Config.Driver, DataSourceName)

	if err != nil {
		log.Panicln(err)
	}

	db.LogMode(config.Common.Mode != "production")

	if Config.Sync == "on" {
		log.Println("正在同步数据库...")

		// Migrate the schema
		db.AutoMigrate(
			new(User),                 // 用户表
			new(OAuth),                // oAuth2 表
			new(Host),                 // 服务器表
			new(HostRecord),           // 服务器许可记录
			new(HostConnectionRecord), // 服务器连接操作记录
			new(Team),                 // 团队表
			new(TeamMember),           // 团队成名表
			new(TeamMemberInvite),     // 团队邀请记录表
		)

		log.Println("数据库同步完成.")
	}

	Db = db
}

func DeleteRowByTable(tableName string, field string, value interface{}) (err error) {
	var (
		tx *gorm.DB
	)

	defer func() {
		if tx != nil {
			if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}
	}()

	tx = Db.Begin()

	raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE %s = '%v'", tableName, field, value)

	if err = tx.Exec(raw).Error; err != nil {
		return
	}

	return
}
