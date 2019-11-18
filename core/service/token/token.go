// Copyright 2019 Axetroy. All rights reserved. MIT license.
package token

import (
	"errors"
	"github.com/axetroy/terminal/core/config"
	"github.com/dgrijalva/jwt-go"
	"log"
)

const (
	Prefix    = "Bearer"
	AuthField = "Authorization"
)

var (
	userSecreteKey  string
	adminSecreteKey string
)

type Claims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

type ClaimsInternal struct {
	Uid string `json:"uid"` // base64 encode
	jwt.StandardClaims
}

func init() {
	userSecreteKey = config.User.Secret
	adminSecreteKey = config.Admin.Secret
	if userSecreteKey == adminSecreteKey {
		log.Panicln(errors.New("用户端的 Token 密钥不能和管理员端的相同，存在安全风险"))
	}
}

func JoinPrefixToken(token string) string {
	return Prefix + " " + token
}
