// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build !windows

package locale

import (
	"os"
	"strings"
)

func getLocaleName() (string, error) {
	name := os.Getenv("LANG")

	// LANG = zh_CN.UTF-8 过滤掉最后的编译方式
	index := strings.LastIndexByte(name, '.')
	if index > 0 {
		name = name[:index]
	}

	return name, nil
}
