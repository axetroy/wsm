// Copyright 2019 Axetroy. All rights reserved. MIT license.
package database

import (
	"fmt"
	"log"

	"github.com/axetroy/terminal/internal/app/config"
	"github.com/axetroy/terminal/internal/app/model"
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
			new(model.User),  // 用户表
			new(model.OAuth), // oAuth2 表
			new(model.Host),  // 服务器表
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
