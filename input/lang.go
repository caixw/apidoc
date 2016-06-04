// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// 所有支持的语言模型定义
//
// NOTE: 应该保持键名为小写
var langs = map[string][]*block{
	// golang
	"go": []*block{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "`", End: "`"},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// c/c++
	"c": []*block{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// php
	"php": []*block{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// ruby
	"ruby": []*block{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `#`},
		// BUG(caixw): 一个单行注释后紧跟前多行注释时，多行注释会被忽略。
		&block{Type: blockTypeMComment, Begin: "\n=begin", End: "\n=end"},
	},
}

// 各语言默认支持的文件扩展名。
//
// NOTE: 应该保持键名、键值均为小写
var langExts = map[string][]string{
	"go":   []string{".go"},
	"c":    []string{".h", ".c", ".cpp", ".cxx"},
	"php":  []string{".php"},
	"ruby": []string{".rb"},
}

// 返回所有支持的语言
func Langs() []string {
	ret := make([]string, 0, len(langs))
	for l := range langs {
		ret = append(ret, l)
	}

	return ret
}

// 检测指定目录下的语言类型。
//
// 检测依据为根据扩展名来做统计，数量最大且被支持的获胜。
func DetectDirLang(dir string) (string, error) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// langsMap 记录每个支持的语言对应的文件数量
	langsMap := make(map[string]int, len(fs))
	for _, f := range fs { // 遍历所有的文件
		if f.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(f.Name()))
		lang := getLangByExt(ext)
		if len(lang) > 0 {
			langsMap[lang]++
		}
	}

	if len(langsMap) == 0 {
		return "", errors.New("该目录下没有支持的语言文件")
	}

	lang := ""
	cnt := 0
	for k, v := range langsMap {
		if v >= cnt {
			lang = k
			cnt = v
		}
	}

	if len(lang) > 0 {
		return lang, nil
	}
	return "", errors.New("该目录下没有支持的语言文件")
}

// 根据扩展名获取其对应的语言名称。
// 若返回空值，则表示没有找到对应的。
func getLangByExt(ext string) string {
	ext = strings.ToLower(ext)
	for lang, exts := range langExts {
		for _, elem := range exts {
			if elem == ext {
				return lang
			}
		}
	}
	return ""
}

// 是否支持该语言
func langIsSupported(lang string) bool {
	// 由测试函数保证 langs 和 langExts 拥有相同的键名。
	_, found := langs[lang]
	return found
}
