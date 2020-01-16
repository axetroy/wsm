// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package token

import (
	"time"

	"github.com/axetroy/wsm/internal/library/util"
	"github.com/dgrijalva/jwt-go"
)

// generate jwt token
func Generate(secret, userId string) (tokenString string, err error) {
	// 生成token
	c := ClaimsInternal{
		util.Base64Encode(userId),
		jwt.StandardClaims{
			Audience:  userId,
			Id:        userId,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(6)).Unix(),
			Issuer:    "user",
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err = token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	tokenString = JoinPrefixToken(tokenString)

	return
}
