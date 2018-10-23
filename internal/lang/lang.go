// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package lang 各类语言解析和管理。
package lang

// 所有支持的语言模型定义
var langs = []*Language{
	&Language{
		Name:   "c#",
		Exts:   []string{".cs"},
		Blocks: cStyle,
	},

	&Language{
		Name:   "c++",
		Exts:   []string{".h", ".c", ".cpp", ".cxx", ".hpp"},
		Blocks: cStyle,
	},

	&Language{
		Name:   "d",
		Exts:   []string{".d"},
		Blocks: cStyle,
	},

	&Language{
		Name: "erlang",
		Exts: []string{".erl", ".hrl"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `%`},
		},
	},

	&Language{
		Name: "go",
		Exts: []string{".go"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeString, Begin: "`", End: "`"},
			&Block{Type: BlockTypeSComment, Begin: `//`},
			&Block{Type: BlockTypeMComment, Begin: `/*`, End: `*/`},
		},
	},

	&Language{
		Name: "groovy",
		Exts: []string{".groovy"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeString, Begin: "'", End: "'", Escape: `\`},
			&Block{Type: BlockTypeString, Begin: "'''", End: "'''", Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `//`},
			&Block{Type: BlockTypeMComment, Begin: `/*`, End: `*/`},
		},
	},

	&Language{
		Name:   "java",
		Exts:   []string{".java"},
		Blocks: cStyle,
	},

	&Language{
		Name: "javascript",
		Exts: []string{".js"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeString, Begin: "'", End: "'", Escape: `\`},
			&Block{Type: BlockTypeString, Begin: "`", End: "`", Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `//`},
			&Block{Type: BlockTypeMComment, Begin: `/*`, End: `*/`},
			// NOTE: js 中若出现 /*abc/.test() 应该是先优先注释的。放最后，优先匹配 // 和 /*
			&Block{Type: BlockTypeString, Begin: "/", End: "/", Escape: `\`}, // 正则表达式
		},
	},

	&Language{
		Name: "pascal",
		Exts: []string{".pas", ".pp"},
		Blocks: []Blocker{
			newPascalStringBlock('\''),
			newPascalStringBlock('"'),
			&Block{Type: BlockTypeMComment, Begin: "{", End: "}"},
			&Block{Type: BlockTypeMComment, Begin: "(*", End: "*)"},
		},
	},

	&Language{
		Name: "perl",
		Exts: []string{".perl", ".prl", ".pl"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeString, Begin: "'", End: "'", Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `#`},
			&Block{Type: BlockTypeMComment, Begin: "\n=pod\n", End: "\n=cut\n"},
		},
	},

	&Language{
		Name: "php",
		Exts: []string{".php"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeString, Begin: "'", End: "'", Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `//`},
			&Block{Type: BlockTypeSComment, Begin: `#`},
			&Block{Type: BlockTypeMComment, Begin: `/*`, End: `*/`},
		},
	},

	&Language{
		Name: "python",
		Exts: []string{".py"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeMComment, Begin: `"""`, End: `"""`},
			&Block{Type: BlockTypeMComment, Begin: "'''", End: `'''`},
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `#`},
		},
	},

	&Language{
		Name: "ruby",
		Exts: []string{".rb"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeString, Begin: "'", End: "'", Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `#`},
			&Block{Type: BlockTypeMComment, Begin: "\n=begin\n", End: "\n=end\n"},
		},
	},

	&Language{
		Name: "rust",
		Exts: []string{".rs"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `///`}, // 需要在 // 之前定义
			&Block{Type: BlockTypeSComment, Begin: `//`},
			&Block{Type: BlockTypeMComment, Begin: `/*`, End: `*/`},
		},
	},

	&Language{
		Name:   "scala",
		Exts:   []string{".scala"},
		Blocks: cStyle,
	},

	&Language{
		Name: "swift",
		Exts: []string{".swift"},
		Blocks: []Blocker{
			&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&Block{Type: BlockTypeSComment, Begin: `//`},
			newSwiftNestMCommentBlock("/*", "*/"),
		},
	},
}

var cStyle = []Blocker{
	&Block{Type: BlockTypeString, Begin: `"`, End: `"`, Escape: `\`},
	&Block{Type: BlockTypeSComment, Begin: `//`},
	&Block{Type: BlockTypeMComment, Begin: `/*`, End: `*/`},
}

// Language 语言模块的定义
type Language struct {
	Name   string
	Blocks []Blocker
	Exts   []string
}

// Get 获取指定语言的定义信息
//
// 若不存在，则返回 nil
func Get(name string) *Language {
	for _, lang := range langs {
		if lang.Name == name {
			return lang
		}
	}

	return nil
}

// GetByExt 根据扩展名获取语言定义信息
//
// 若不存在，则返回 nil
func GetByExt(ext string) *Language {
	for _, lang := range langs {
		for _, e := range lang.Exts {
			if e == ext {
				return lang
			}
		}
	}

	return nil
}

// Langs 返回所有支持的语言
func Langs() []string {
	ret := make([]string, 0, len(langs))
	for _, l := range langs {
		ret = append(ret, l.Name)
	}

	return ret
}
