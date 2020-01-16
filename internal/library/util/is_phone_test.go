// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package util_test

import (
	"testing"

	"github.com/axetroy/wsm/internal/library/util"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Input  string
	Expect bool
}

func TestIsPhone(t *testing.T) {
	tests := []testCase{
		{
			Input:  "13333333333",
			Expect: true,
		},
		{
			Input:  "133333333331",
			Expect: false,
		},
		{
			Input:  "03333333333",
			Expect: false,
		},
	}

	for _, input := range tests {
		assert.Equal(t, input.Expect, util.IsPhone(input.Input))
	}
}
