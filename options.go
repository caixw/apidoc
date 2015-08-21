// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caixw/apidoc/core"
)

// 所有需要传递给scanner包的参数集合。
type Options struct {
	SrcDir    string
	Recursive bool
	Type      string
	Exts      []string
	Groups    []string
}

// 通过参数作一些初始化工作。
// 从Options实例中获取真正需要的参数。
func getArgs(opt *Options) (core.ScanFunc, []string, error) {
	dir := opt.SrcDir + string(os.PathSeparator)

	exts := make([]string, 0, len(opt.Exts))
	for _, ext := range opt.Exts {
		if len(ext) == 0 {
			continue
		}

		if ext[0] != '.' {
			ext = "." + ext
		}
		exts = append(exts, ext)
	}

	typ := opt.Type
	// 若没有指定Type，则根据exts和当前目录下的文件检测来确定其值
	if len(opt.Type) == 0 {
		var err error
		if len(exts) == 0 {
			typ, err = detectDirLangType(dir)
		} else {
			typ, err = detectLangType(exts)
		}

		if err != nil {
			return nil, nil, err
		}
	}

	typ = strings.ToLower(typ)
	lang, found := langs[typ]
	if !found {
		return nil, nil, fmt.Errorf("暂不支持该类型[%v]的语言", typ)
	}
	if len(exts) == 0 {
		exts = lang.exts
	}

	paths, err := recursivePath(dir, opt.Recursive, exts...)
	if err != nil {
		return nil, nil, err
	}

	return lang.scan, paths, nil
}

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
