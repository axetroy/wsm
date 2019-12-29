// Copyright 2019 Axetroy. All rights reserved. Apache License 2.0.
package redis

import (
	"github.com/axetroy/wsm/internal/app/config"
	"github.com/go-redis/redis"
)

var (
	Client          *redis.Client // 默认的redis存储
	ClientOAuthCode *redis.Client // 存储 oAuth2 对应的激活码
	Config          = config.Redis
)

func init() {
	var (
		addr     = Config.Host + ":" + Config.Port
		password = Config.Password
	)

	// 初始化DB连接
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	ClientOAuthCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       5,
	})

}
