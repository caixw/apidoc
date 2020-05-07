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
			newCStyleString(),
			newSingleComment("%"),
		},
	},

	{
		DisplayName: "Go",
		Name:        "go",
		Exts:        []string{".go"},
		Blocks: []Blocker{
			newCStyleString(),
			newString("`", "`", ""),
			newCStyleChar(),
			newCStyleSingleComment(),
			newCStyleMultipleComment(),
		},
	},

	{
		DisplayName: "Groovy",
		Name:        "groovy",
		Exts:        []string{".groovy"},
		Blocks: []Blocker{
			newCStyleString(),
			newString("'", "'", `\`),
			newString("'''", "'''", `\`),
			newCStyleSingleComment(),
			newCStyleMultipleComment(),
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
			newCStyleString(),
			newString("'", "'", `\`),
			newString("`", "`", `\`),
			newCStyleSingleComment(),
			newCStyleMultipleComment(),
			// NOTE: js 中若出现 /*abc/.test() 应该是先优先注释的。放最后，优先匹配 // 和 /*
			newString("/", "/", `\`),
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
			newMultipleComment("{", "}", ""),
			newMultipleComment("(*", "*)", "*"),
		},
	},

	{
		DisplayName: "Perl",
		Name:        "perl",
		Exts:        []string{".perl", ".prl", ".pl"},
		Blocks: []Blocker{
			newCStyleString(),
			newString("'", "'", `\`),
			newSingleComment("#"),
			newRubyMultipleComment("=pod", "=cut", ""),
		},
	},

	{
		DisplayName: "PHP",
		Name:        "php",
		Exts:        []string{".php"},
		Blocks: []Blocker{
			newCStyleString(),
			newString("'", "'", `\`),
			newPHPDocBlock(),
			newCStyleSingleComment(),
			newSingleComment("#"),
			newCStyleMultipleComment(),
		},
	},

	{
		DisplayName: "Python",
		Name:        "python",
		Exts:        []string{".py"},
		Blocks: []Blocker{
			newMultipleComment(`"""`, `"""`, ""),
			newMultipleComment("'''", "'''", ""),
			newCStyleString(),
			newSingleComment(`#`),
		},
	},

	{
		DisplayName: "Ruby",
		Name:        "ruby",
		Exts:        []string{".rb"},
		Blocks: []Blocker{
			newCStyleString(),
			newString("'", "'", `\`),
			newSingleComment(`#`),
			newRubyMultipleComment("=begin", "=end", ""),
		},
	},

	{
		DisplayName: "Rust",
		Name:        "rust",
		Exts:        []string{".rs"},
		Blocks:      cStyle,
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
			newString(`"""`, `"""`, `\`),
			newString(`#"`, `"#`, ""),
			newCStyleString(),
			newCStyleChar(),
			newCStyleSingleComment(),
			newSwiftNestMCommentBlock("/*", "*/", "*"),
		},
	},
}

var cStyle = []Blocker{
	newCStyleString(),
	newCStyleChar(),
	newSingleComment("///"), // 需要在 // 之前定义
	newCStyleSingleComment(),
	newCStyleMultipleComment(),
}

// 处理 "XXX\""
func newCStyleString() Blocker {
	return newString(`"`, `"`, `\`)
}

// 处理 '"'
func newCStyleChar() Blocker {
	return newString(`'`, `'`, "")
}

// 处理 // xxx
func newCStyleSingleComment() Blocker {
	return newSingleComment(`//`)
}

// 处理 /* */
func newCStyleMultipleComment() Blocker {
	return newMultipleComment(`/*`, `*/`, "*")
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
