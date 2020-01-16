// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package config

import (
	"github.com/axetroy/wsm/internal/library/dotenv"
)

type TLS struct {
	Cert string `json:"cert"` // 证书文件
	Key  string `json:"key"`  // Key 文件
}

type user struct {
	Domain string `json:"domain"` // 用户端 API 绑定的域名, 例如 https://example.com
	Port   string `json:"port"`   // 用户端 API 监听的端口
	Secret string `json:"secret"` // 用户端密钥，用于加密/解密 token
	TLS    *TLS   `json:"tls"`
}

var User user

func init() {
	User.Port = dotenv.GetByDefault("PORT", "8080")
	User.Secret = dotenv.GetByDefault("SECRET_KEY", "user")

	TlsCert := dotenv.GetByDefault("TLS_CERT", "")
	TlsKey := dotenv.GetByDefault("TLS_KEY", "")

	if TlsCert != "" && TlsKey != "" {
		User.TLS = &TLS{
			Cert: TlsCert,
			Key:  TlsKey,
		}
	}
}
