// SPDX-License-Identifier: MIT

// Package lang 管理各类语言提取注释代码块规则的定义
package lang

import "fmt"

// 所有支持的语言模型定义
var langs = []*Language{
	{
		DisplayName: "C#",
		Name:        "c#",
		Exts:        []string{".cs"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "C/C++",
		Name:        "c++",
		Exts:        []string{".h", ".c", ".cpp", ".cxx", ".hpp"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "D",
		Name:        "d",
		Exts:        []string{".d"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "Erlang",
		Name:        "erlang",
		Exts:        []string{".erl", ".hrl"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeSComment, Begin: `%`},
		},
	},

	{
		DisplayName: "Go",
		Name:        "go",
		Exts:        []string{".go"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeString, Begin: "`", End: "`"},
			&block{Type: blockTypeString, Begin: `'`, End: `'`}, // 处理 '"‘ 的内容
			&block{Type: blockTypeSComment, Begin: `//`},
			&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`, Escape: "*"},
		},
	},

	{
		DisplayName: "Groovy",
		Name:        "groovy",
		Exts:        []string{".groovy"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
			&block{Type: blockTypeString, Begin: "'''", End: "'''", Escape: `\`},
			&block{Type: blockTypeSComment, Begin: `//`},
			&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`, Escape: "*"},
		},
	},

	{
		DisplayName: "Java",
		Name:        "java",
		Exts:        []string{".java"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "JavaScript",
		Name:        "javascript",
		Exts:        []string{".js"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
			&block{Type: blockTypeString, Begin: "`", End: "`", Escape: `\`},
			&block{Type: blockTypeSComment, Begin: `//`},
			&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`, Escape: "*"},
			// NOTE: js 中若出现 /*abc/.test() 应该是先优先注释的。放最后，优先匹配 // 和 /*
			&block{Type: blockTypeString, Begin: "/", End: "/", Escape: `\`}, // 正则表达式
		},
	},

	{
		DisplayName: "Kotlin",
		Name:        "kotlin",
		Exts:        []string{".kt"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "Pascal/Delphi",
		Name:        "pascal",
		Exts:        []string{".pas", ".pp"},
		Blocks: []Blocker{
			newPascalStringBlock('\''),
			newPascalStringBlock('"'),
			&block{Type: blockTypeMComment, Begin: "{", End: "}"},
			&block{Type: blockTypeMComment, Begin: "(*", End: "*)", Escape: "*"},
		},
	},

	{
		DisplayName: "Perl",
		Name:        "perl",
		Exts:        []string{".perl", ".prl", ".pl"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
			&block{Type: blockTypeSComment, Begin: `#`},
			&block{Type: blockTypeMComment, Begin: "\n=pod\n", End: "\n=cut\n"},
		},
	},

	{
		DisplayName: "PHP",
		Name:        "php",
		Exts:        []string{".php"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
			newPHPDocBlock(),
			&block{Type: blockTypeSComment, Begin: `//`},
			&block{Type: blockTypeSComment, Begin: `#`},
			&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`, Escape: "*"},
		},
	},

	{
		DisplayName: "Python",
		Name:        "python",
		Exts:        []string{".py"},
		Blocks: []Blocker{
			&block{Type: blockTypeMComment, Begin: `"""`, End: `"""`},
			&block{Type: blockTypeMComment, Begin: "'''", End: `'''`},
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeSComment, Begin: `#`},
		},
	},

	{
		DisplayName: "Ruby",
		Name:        "ruby",
		Exts:        []string{".rb"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeString, Begin: "'", End: "'", Escape: `\`},
			&block{Type: blockTypeSComment, Begin: `#`},
			&block{Type: blockTypeMComment, Begin: "\n=begin\n", End: "\n=end\n"},
		},
	},

	{
		DisplayName: "Rust",
		Name:        "rust",
		Exts:        []string{".rs"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeString, Begin: `'`, End: `'`}, // 处理 '"‘ 的内容
			&block{Type: blockTypeSComment, Begin: `///`},       // 需要在 // 之前定义
			&block{Type: blockTypeSComment, Begin: `//`},
			&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`, Escape: "*"},
		},
	},

	{
		DisplayName: "Scala",
		Name:        "scala",
		Exts:        []string{".scala"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "Swift",
		Name:        "swift",
		Exts:        []string{".swift"},
		Blocks: []Blocker{
			&block{Type: blockTypeString, Begin: `"""`, End: `"""`, Escape: `\`},
			&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
			&block{Type: blockTypeString, Begin: `#"`, End: `"#`},
			&block{Type: blockTypeString, Begin: `'`, End: `'`}, // 处理 '"‘ 的内容
			&block{Type: blockTypeSComment, Begin: `//`},
			newSwiftNestMCommentBlock("/*", "*/", "*"),
		},
	},
}

var cStyle = []Blocker{
	&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
	&block{Type: blockTypeString, Begin: `'`, End: `'`}, // 处理 '"‘ 的内容
	&block{Type: blockTypeSComment, Begin: `///`},       // 需要在 // 之前定义
	&block{Type: blockTypeSComment, Begin: `//`},
	&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`, Escape: "*"},
}

// Language 语言模块的定义
type Language struct {
	DisplayName string    // 显示友好的名称
	Name        string    // 语言唯一名称，一律小写
	Blocks      []Blocker // 注释块的解析规则定义
	Exts        []string  // 扩展名列表，必须以 . 开头且小写
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
// ext 必须以 . 作为开头
// 若不存在，则返回 nil
func GetByExt(ext string) *Language {
	if len(ext) == 0 || ext[0] != '.' {
		panic(fmt.Sprintf("参数 ext 的值 [%s] 不能为空，且必须以 . 作为开头", ext))
	}

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
func Langs() []*Language {
	return langs
}
