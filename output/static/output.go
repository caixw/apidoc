// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package static

import (
	"io/ioutil"
	"os"
)

// 将所有的静态文件输出到该目录下。
func Output(dir string) error {
	for path, content := range files {
		path = dir + string(os.PathSeparator) + path
		if err := ioutil.WriteFile(path, content, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
