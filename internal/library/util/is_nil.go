// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package util

import "reflect"

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
