// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package util

import (
	"strconv"
)

// 统一的货币金额格式化，保留 8 位小数
func FloatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', 8, 64)
}
