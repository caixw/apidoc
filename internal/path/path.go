// SPDX-License-Identifier: MIT

package path

import (
	"os"
	"path/filepath"
	"strings"
)

// Abs 获取 path 的绝对路径
//
// 如果 path 是相对路径的，则将其设置为相对于 wd 的路径
func Abs(path, wd string) (p string, err error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	if !isBeginHome(path) {
		path = filepath.Join(wd, path)
	}

	// ~ 路开头的相对路径，需要将其定位到 HOME 目录之下
	if isBeginHome(path) {
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

func isBeginHome(path string) bool {
	return strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\")
}

// Rel 尽可能地返回 path 相对于 wd 的路径，如果不存在相对关系，则原因返回 path。
func Rel(path, wd string) string {
	p, err := filepath.Rel(wd, path)
	if err != nil { // 不能转换不算错误，直接返回原值
		return path
	}
	return p
}
