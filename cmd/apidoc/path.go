// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"path/filepath"
	"strings"
)

// abs 获取 path 的绝对路径
//
// 如果 path 是相对路径的，则将其设置为相对于 wd 的路径
//
// wd 表示工作目录，当 path 不是绝对路径和 ~ 开头时，表示相对于此目录；
// path 表示需要处理的路径。
func abs(path, wd string) (p string, err error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	if !strings.HasPrefix(path, "~/") {
		path = filepath.Join(wd, path)
	}

	// 非 ~ 路开头的相对路径，需要将其定位到 wd 目录之下
	if strings.HasPrefix(path, "~/") {
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(dir, path[2:])
	}

	if !filepath.IsAbs(path) {
		if path, err = filepath.Abs(path); err != nil {
			return "", err
		}
	}

	return filepath.Clean(path), nil
}
