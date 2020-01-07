// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package redis

import (
	"github.com/axetroy/wsm/internal/app/config"
	"github.com/go-redis/redis"
)

var (
	ClientConnection *redis.Client // 存储 对应的 SSH 连接∏
	ClientOAuthCode  *redis.Client // 存储 oAuth2 对应的激活码
	Config           = config.Redis
)

func init() {
	var (
		addr     = Config.Host + ":" + Config.Port
		password = Config.Password
	)

	// 初始化DB连接
	ClientConnection = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1,
	})

	ClientOAuthCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       5,
	})

}
