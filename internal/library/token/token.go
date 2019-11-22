// Copyright 2019 Axetroy. All rights reserved. MIT license.
package token

import (
	"github.com/dgrijalva/jwt-go"
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

func JoinPrefixToken(token string) string {
	return Prefix + " " + token
}
