// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package fs

import "os"

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
