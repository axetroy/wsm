// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/terminal/core/service/dotenv"
)

type database struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Driver       string `json:"driver"`
	DatabaseName string `json:"database_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Sync         string `json:"sync"`
}

var Database database

func init() {
	Database.Driver = dotenv.GetByDefault("DB_DRIVER", "postgres")
	Database.Host = dotenv.GetByDefault("DB_HOST", "localhost")
	Database.Port = dotenv.GetByDefault("DB_PORT", "54321")
	Database.DatabaseName = dotenv.GetByDefault("DB_NAME", "terminal")
	Database.Username = dotenv.GetByDefault("DB_USERNAME", "terminal")
	Database.Password = dotenv.GetByDefault("DB_PASSWORD", "terminal")
	Database.Sync = dotenv.GetByDefault("DB_SYNC", "off")
}
