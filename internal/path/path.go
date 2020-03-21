// SPDX-License-Identifier: MIT

// Package path 提供一些文件相关的操作
package path

import (
	"path/filepath"
	"runtime"
)

// CurrPath 获取相当于调用者所在目录的路径列表，相当于 PHP 的 __DIR__ + "/" + path
func CurrPath(path string) string {
	_, fi, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(fi), path)
}
