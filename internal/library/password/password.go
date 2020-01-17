// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	prefix = "ped&13()%"
	suffix = "d;'1^3@!#"
)

// 生成密码
func Generate(plaintext string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(prefix+plaintext+suffix), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// 校验密码
func Verify(plaintext, hash string) bool {
	// 错误密码验证
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(prefix+plaintext+suffix))

	if err != nil {
		return false
	} else {
		return true
	}
}
