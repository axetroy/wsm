// Copyright 2019 Axetroy. All rights reserved. MIT license.
package database

import (
	"fmt"
	"log"

	"github.com/axetroy/terminal/core/config"
	"github.com/axetroy/terminal/core/model"
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
			new(model.Role),  // 角色表 - RBAC
			new(model.OAuth), // oAuth2 表
		)

		log.Println("数据库同步完成.")
	}

	Db = db

	defaultRole := model.Role{Name: model.DefaultUser.Name}

	// 确保有默认的角色
	if err := db.First(&defaultRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&model.Role{
				Name:        model.DefaultUser.Name,
				Description: model.DefaultUser.Description,
				Accession:   model.DefaultUser.AccessionArray(),
				BuildIn:     true,
			}).Error
		} else {
			log.Panicln(err)
		}
	} else {
		// 如果角色已存在，则同步角色的权限
		if err := db.Model(&defaultRole).Update(&model.Role{
			Accession: model.DefaultUser.AccessionArray(),
		}).Error; err != nil {
			log.Panicln(err)
		}
	}

}

func DeleteRowByTable(tableName string, field string, value interface{}) {
	var (
		err error
		tx  *gorm.DB
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
}
