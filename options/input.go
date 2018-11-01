// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package options

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/caixw/apidoc/internal/lang"
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
func Detect(dir string, recursive bool) ([]*Input, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	exts, err := detectExts(dir, recursive)
	if err != nil {
		return nil, err
	}

	langs := detectLanguage(exts)

	inputs := make([]*Input, 0, len(langs))
	for _, lang := range langs {
		inputs = append(inputs, &Input{
			Lang:      lang.Name,
			Dir:       dir,
			Exts:      lang.Exts,
			Recursive: recursive,
		})
	}

	return inputs, nil
}

type language struct {
	lang.Language
	count int
}

// 根据 exts 计算每个语言对应的文件数量，并按倒序返回
//
// exts 参数为从 detectExts 中获取的返回值
func detectLanguage(exts map[string]int) []*language {
	langs := make([]*language, 0, len(exts))

	for ext, count := range exts {
		l := lang.GetByExt(ext)
		if l == nil {
			continue
		}

		found := false
		for _, item := range langs {
			if item.Name == l.Name {
				item.count += count
				found = true
				break
			}
		}
		if !found {
			langs = append(langs, &language{
				count:    count,
				Language: *l,
			})
		}
	} // end for

	sort.SliceStable(langs, func(i, j int) bool {
		return langs[i].count > langs[j].count
	})

	return langs
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
			if len(ext) > 0 {
				exts[ext]++
			}
		}

		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	return exts, nil
}
