// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/caixw/apidoc/locale"
)

// 所有支持的语言模型定义
//
// NOTE: 应该保持键名为非大写，按字母顺序排列，方便查找。
// langs 应该和 langExts 保持一一对应关系。
var langs = map[string][]blocker{
	// C#
	"c#": cStyle,

	// c/c++
	"c++": cStyle,

	// d
	"d": cStyle,

	// erlang
	"erlang": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `%`},
	},

	// golang
	"go": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "`", End: "`"},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// groovy
	"groovy": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeString, Begin: "'''", End: "'''", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// java
	"java": cStyle,

	// javascript
	"javascript": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
		// NOTE: js 中若出现 /*abc/.test() 应该是先优先注释的。放最后，优先匹配 // 和 /*
		&block{Type: blockTypeString, Begin: "/", End: "/", Escape: `\`}, // 正则表达式
	},

	// pascal/delphi
	"pascal": []blocker{
		newPascalStringBlock('\''),
		newPascalStringBlock('"'),
		&block{Type: blockTypeMComment, Begin: "{", End: "}"},
		&block{Type: blockTypeMComment, Begin: "(*", End: "*)"},
	},

	// perl
	"perl": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `#`},
		&block{Type: blockTypeMComment, Begin: "\n=pod\n", End: "\n=cut\n"},
	},

	// python
	"python": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `#`},
	},

	// php
	"php": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// ruby
	"ruby": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `#`},
		&block{Type: blockTypeMComment, Begin: "\n=begin\n", End: "\n=end\n"},
	},

	// rust
	"rust": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `///`}, // 需要在 // 之前定义
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// scala
	"scala": cStyle,

	// swift
	"swift": []blocker{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		newSwiftNestMCommentBlock("/*", "*/"),
	},
}

var cStyle = []blocker{
	&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
	&block{Type: blockTypeSComment, Begin: `//`},
	&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
}

// 各语言默认支持的文件扩展名。
//
// NOTE: 键名与 langs 中的键一一对应。
// 键值为可能的文件扩展名列表，应该使用非大写状态。
// `go test` 会作大量名称是否规范的检测。
var langExts = map[string][]string{
	"c#":         []string{".cs"},
	"c++":        []string{".h", ".c", ".cpp", ".cxx", "hpp"},
	"d":          []string{".d"},
	"erlang":     []string{".erl", "hrl"},
	"go":         []string{".go"},
	"groovy":     []string{".groovy"},
	"java":       []string{".java"},
	"javascript": []string{".js"},
	"pascal":     []string{".pas", ".pp"},
	"perl":       []string{".perl", ".prl", ".pl"},
	"php":        []string{".php"},
	"python":     []string{".py"},
	"ruby":       []string{".rb"},
	"rust":       []string{".rs"},
	"scala":      []string{".scala"},
	"swift":      []string{".swift"},
}

// Langs 返回所有支持的语言
func Langs() []string {
	ret := make([]string, 0, len(langs))
	for l := range langs {
		ret = append(ret, l)
	}

	return ret
}

// DetectDirLang 检测指定目录下的语言类型。
//
// 检测依据为根据扩展名来做统计，数量最大且被支持的获胜。
// 不会分析子目录。
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
		return "", errors.New(locale.Sprintf(locale.ErrNotFoundSupportedLang))
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
	return "", errors.New(locale.Sprintf(locale.ErrNotFoundSupportedLang))
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
