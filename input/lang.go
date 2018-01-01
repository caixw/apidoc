// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"sort"
	"strings"
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
	"erlang": {
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `%`},
	},

	// go
	"go": {
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "`", End: "`"},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// groovy
	"groovy": {
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeString, Begin: "'''", End: "'''", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// java
	"java": cStyle,

	// javascript
	"javascript": {
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeString, Begin: "`", End: "`", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
		// NOTE: js 中若出现 /*abc/.test() 应该是先优先注释的。放最后，优先匹配 // 和 /*
		&block{Type: blockTypeString, Begin: "/", End: "/", Escape: `\`}, // 正则表达式
	},

	// pascal/delphi
	"pascal": {
		newPascalStringBlock('\''),
		newPascalStringBlock('"'),
		&block{Type: blockTypeMComment, Begin: "{", End: "}"},
		&block{Type: blockTypeMComment, Begin: "(*", End: "*)"},
	},

	// perl
	"perl": {
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `#`},
		&block{Type: blockTypeMComment, Begin: "\n=pod\n", End: "\n=cut\n"},
	},

	// python
	"python": {
		&block{Type: blockTypeMComment, Begin: `"""`, End: `"""`},
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `#`},
	},

	// php
	"php": {
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// ruby
	"ruby": {
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `#`},
		&block{Type: blockTypeMComment, Begin: "\n=begin\n", End: "\n=end\n"},
	},

	// rust
	"rust": {
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `///`}, // 需要在 // 之前定义
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// scala
	"scala": cStyle,

	// swift
	"swift": {
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
	"c#":         {".cs"},
	"c++":        {".h", ".c", ".cpp", ".cxx", ".hpp"},
	"d":          {".d"},
	"erlang":     {".erl", ".hrl"},
	"go":         {".go"},
	"groovy":     {".groovy"},
	"java":       {".java"},
	"javascript": {".js"},
	"pascal":     {".pas", ".pp"},
	"perl":       {".perl", ".prl", ".pl"},
	"php":        {".php"},
	"python":     {".py"},
	"ruby":       {".rb"},
	"rust":       {".rs"},
	"scala":      {".scala"},
	"swift":      {".swift"},
}

// Languages 返回所有支持的语言
func Languages() []string {
	ret := make([]string, 0, len(langs))
	for l := range langs {
		ret = append(ret, l)
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i] > ret[j]
	})

	return ret
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
