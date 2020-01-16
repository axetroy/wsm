// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package validator

import (
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/library/util"
	"regexp"
)

var (
	usernameReg = regexp.MustCompile("^[\\w\\-]+$")
)

func ValidateStruct(any interface{}) error {
	if isValid, err := govalidator.ValidateStruct(any); err != nil {
		return exception.New(err.Error(), exception.InvalidParams.Code())
	} else if !isValid {
		return exception.InvalidParams
	}
	return nil
}

func IsEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func IsPhone(phone string) bool {
	return util.IsPhone(phone)
}

func IsValidUsername(username string) bool {
	return usernameReg.MatchString(username)
}

func ValidatePhone(phone string) error {
	if !IsPhone(phone) {
		return exception.InvalidFormat
	} else {
		return nil
	}
}

func ValidateEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return exception.InvalidFormat
	} else {
		return nil
	}
}

func ValidateUsername(username string) error {
	if !IsValidUsername(username) {
		return exception.InvalidFormat
	}
	return nil
}
