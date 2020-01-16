// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package tester

import (
	"encoding/json"
	"errors"

	"github.com/axetroy/wsm/internal/library/util"
)

func Decode(source interface{}, dest interface{}) (err error) {
	if !util.IsPoint(dest) {
		err = errors.New("decode: dest expect a point")
		return
	}

	var b []byte
	if b, err = json.Marshal(source); err != nil {
		return
	}

	if err = json.Unmarshal(b, dest); err != nil {
		return
	}
	return
}
