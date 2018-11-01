// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package options

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caixw/apidoc/internal/lang"
	"github.com/caixw/apidoc/internal/locale"
)

// Input 指定输入内容的相关信息。
type Input struct {
	Lang      string   `yaml:"lang"`               // 输入的目标语言
	Dir       string   `yaml:"dir"`                // 源代码目录，建议使用绝对路径
	Exts      []string `yaml:"exts,omitempty"`     // 需要扫描的文件扩展名，若未指定，则使用默认值
	Recursive bool     `yaml:"recursive"`          // 是否查找 Dir 的子目录
	Encoding  string   `yaml:"encoding,omitempty"` // 文件的编码，为空表示 utf-8
}

// Detect 检测指定目录下的内容，并为其生成一个合适的 Input 实例。
//
// 检测依据为根据扩展名来做统计，数量最大且被支持的获胜。
func Detect(dir string, recursive bool) (*Input, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	exts, err := detectExts(dir, recursive)
	if err != nil {
		return nil, err
	}

	// 删除不支持的扩展名
	for ext := range exts {
		language := lang.GetByExt(ext)
		if language == nil {
			delete(exts, ext)
		}
	}

	if len(exts) == 0 {
		return nil, locale.Errorf(locale.ErrNotFoundSupportedLang)
	}

	ext := ""
	cnt := 0
	for k, v := range exts {
		if v >= cnt {
			ext = k
			cnt = v
		}
	}
	if len(ext) == 0 {
		return nil, locale.Errorf(locale.ErrNotFoundSupportedLang)
	}

	language := lang.GetByExt(ext)
	if language == nil {
		return nil, locale.Errorf(locale.ErrNotFoundSupportedLang)
	}

	return &Input{
		Lang:      language.Name,
		Dir:       dir,
		Exts:      language.Exts,
		Recursive: recursive,
	}, nil
}

// 返回 dir 目录下文件类型及对应的文件数量的一个集合。
// recursive 表示是否查找子目录。
func detectExts(dir string, recursive bool) (map[string]int, error) {
	exts := map[string]int{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if !recursive && dir != path {
				return filepath.SkipDir
			}
		} else {
			ext := strings.ToLower(filepath.Ext(path))
			exts[ext]++
		}

		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	return exts, nil
}
