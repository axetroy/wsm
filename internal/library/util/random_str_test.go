// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package util_test

import (
	"regexp"
	"testing"

	"github.com/axetroy/wsm/internal/library/util"
	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	assert.Len(t, util.RandomString(1), 1)
	assert.Len(t, util.RandomString(2), 2)
	assert.Len(t, util.RandomString(3), 3)
	assert.Len(t, util.RandomString(4), 4)
	assert.Len(t, util.RandomString(8), 8)
	assert.Len(t, util.RandomString(16), 16)
	assert.IsType(t, "string", util.RandomString(16))
}

func TestRandomNumeric(t *testing.T) {
	assert.Len(t, util.RandomNumeric(1), 1)
	assert.Len(t, util.RandomNumeric(2), 2)
	assert.Len(t, util.RandomNumeric(3), 3)
	assert.Len(t, util.RandomNumeric(4), 4)
	assert.Len(t, util.RandomNumeric(8), 8)
	assert.Len(t, util.RandomNumeric(16), 16)
	assert.IsType(t, "string", util.RandomNumeric(16))
	assert.True(t, regexp.MustCompile("^\\d+$").MatchString(util.RandomNumeric(32)))
}
