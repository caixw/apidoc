// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 默认模板
package static

//go:generate go run make.go

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Output 将所有的静态文件输出到该目录下。
func Output(dir string) error {
	for path, content := range assets {
		path = filepath.Join(dir, path)
		if err := ioutil.WriteFile(path, content, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
