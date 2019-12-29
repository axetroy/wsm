// Copyright 2019 Axetroy. All rights reserved. MIT license.
package token

import (
	"time"

	"github.com/axetroy/wsm/internal/library/util"
	"github.com/dgrijalva/jwt-go"
)

// generate jwt token
func Generate(userId string, isAdmin bool) (tokenString string, err error) {
	var (
		issuer string
		key    string
	)

	if isAdmin {
		issuer = "admin"
		key = adminSecreteKey
	} else {
		issuer = "user"
		key = userSecreteKey
	}

	// 生成token
	c := ClaimsInternal{
		util.Base64Encode(userId),
		jwt.StandardClaims{
			Audience:  userId,
			Id:        userId,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(6)).Unix(),
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err = token.SignedString([]byte(key))

	tokenString = Prefix + " " + tokenString

	return
}
