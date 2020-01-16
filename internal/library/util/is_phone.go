// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package util

import "regexp"

var (
	phoneReg = regexp.MustCompile("^1\\d{10}$")
)

func IsPhone(phoneNumber string) bool {
	return phoneReg.MatchString(phoneNumber)
}
