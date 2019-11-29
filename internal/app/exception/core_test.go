// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception_test

import (
	"fmt"
	"testing"

	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	r := exception.New("test", 123)

	assert.Equal(t, 123, r.Code())
	assert.Equal(t, fmt.Sprintf("test"), r.Error())
	assert.Equal(t, 123, r.Code())
}
