// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package util_test

import (
	"testing"

	"github.com/axetroy/wsm/internal/library/util"
)

func TestMD5(t *testing.T) {
	if util.MD5("1") != "c4ca4238a0b923820dcc509a6f75849b" {
		t.Error("1的MD5值不对")
		return
	}

	if util.MD5("123") != "202cb962ac59075b964b07152d234b70" {
		t.Error("123的MD5值不对")
		return
	}

	if util.MD5("abc") != "900150983cd24fb0d6963f7d28e17f72" {
		t.Error("abc的MD5值不对")
		return
	}
}
