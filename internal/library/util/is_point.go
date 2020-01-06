// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package util

import "reflect"

func IsPoint(i interface{}) bool {
	vi := reflect.ValueOf(i)
	return vi.Kind() == reflect.Ptr
}
