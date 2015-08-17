// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 从扩展名检测其所属的语言名称。
// 以第一个匹配extsIndex的文件扩展名为准。
func detectLangType(exts []string) (string, error) {
	for _, ext := range exts {
		if lang, found := extsIndex[ext]; found {
			return lang, nil
		}
	}
	return "", fmt.Errorf("无法找到与这些扩展名[%v]相匹配的代码扫描函数", exts)
}

// 检测目录下的文件类型。
// 以第一个匹配extsIndex的文件扩展名为准。
func detectDirLangType(dir string) (string, error) {
	var lang string

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() || len(lang) > 0 {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		lang, _ = extsIndex[ext]
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return "", err
	}

	if len(lang) == 0 {
		return lang, fmt.Errorf("无法检测到[%v]目录下的文件类型", dir)
	}

	return lang, nil
}

// 根据recursive值确定是否递归查找paths每个目录下的子目录。
func recursivePath(dir string, recursive bool, exts ...string) ([]string, error) {
	paths := []string{}
	dir += string(os.PathSeparator)

	extIsEnabled := func(ext string) bool {
		for _, v := range exts {
			if ext == v {
				return true
			}
		}
		return false
	}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() && !recursive && path != dir {
			return filepath.SkipDir
		} else if extIsEnabled(filepath.Ext(path)) {
			paths = append(paths, path)
		}
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	return paths, nil
}
