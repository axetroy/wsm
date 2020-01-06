// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package exception_test

import (
	"fmt"
	"testing"

	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	r := exception.New("test", 123)

	assert.Equal(t, 123, r.Code())
	assert.Equal(t, fmt.Sprintf("test"), r.Error())
	assert.Equal(t, 123, r.Code())
}
