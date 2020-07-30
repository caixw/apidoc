// SPDX-License-Identifier: MIT

// Package lang 管理各类语言提取注释代码块规则的定义
package lang

import "fmt"

// 所有支持的语言模型定义
var langs = []*Language{
	{
		DisplayName: "C#",
		ID:          "c#",
		Exts:        []string{".cs"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "C/C++",
		ID:          "c++",
		Exts:        []string{".h", ".c", ".cpp", ".cxx", ".hpp"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "D",
		ID:          "d",
		Exts:        []string{".d"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "Dart",
		ID:          "dart",
		Exts:        []string{".dart"},
		Blocks: []Blocker{
			newCStyleString(),
			newString("'", "'", `\`),
			newString("r'", "'", ``),
			newString("r\"", "\"", ``),
			newString("'''", "'''", ``),
			newString(`"""`, `"""`, ``),
			newSingleComment(`///`),
			newSwiftNestMCommentBlock("/**", "*/", "*"),
			newCStyleSingleComment(),
			newCStyleMultipleComment(),
		},
	},

	{
		DisplayName: "Erlang",
		ID:          "erlang",
		Exts:        []string{".erl", ".hrl"},
		Blocks: []Blocker{
			newCStyleString(),
			newSingleComment("%"),
		},
	},

	{
		DisplayName: "Go",
		ID:          "go",
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
		ID:          "groovy",
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
		ID:          "java",
		Exts:        []string{".java"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "JavaScript",
		ID:          "javascript",
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
		ID:          "kotlin",
		Exts:        []string{".kt"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "Lisp/Clojure",
		ID:          "lisp",
		// .ss,.scm ==> scheme
		// .clj     ==> Coljure
		Exts: []string{".lisp", ".lsp", ".l", ".ss", ".scm", ".clj"},
		Blocks: []Blocker{
			newString(`"`, `"`, `\`), // #"" 为正则表达式
			newSingleComment(";;;;"),
			newSingleComment(";;;"),
			newSingleComment(";;"),
			newSingleComment(";"),
		},
	},

	{
		DisplayName: "Lua",
		ID:          "lua",
		Exts:        []string{".lua"},
		Blocks: []Blocker{
			newString("'", "'", `\`),
			newString("\"", "\"", `\`),
			newString("[[", "]]", ``),
			newSingleComment("--"),
			newMultipleComment("--[[", "]]", "-="),
			newMultipleComment("--[====[", "]====]", "-="),
		},
	},

	{
		DisplayName: "Pascal/Delphi",
		ID:          "pascal",
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
		ID:          "perl",
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
		ID:          "php",
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
		ID:          "python",
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
		ID:          "ruby",
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
		ID:          "rust",
		Exts:        []string{".rs"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "Scala",
		ID:          "scala",
		Exts:        []string{".scala"},
		Blocks:      cStyle,
	},

	{
		DisplayName: "Swift",
		ID:          "swift",
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

	{
		DisplayName: "Typescript",
		ID:          "typescript",
		Exts:        []string{".ts"},
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
		DisplayName: "Zig",
		ID:          "zig",
		Exts:        []string{".zig"},
		Blocks: []Blocker{
			newCStyleString(),
			newCStyleChar(),
			newSingleComment("///"), // 需要在 // 之前定义
			newSingleComment("//!"), // 需要在 // 之前定义
			newCStyleSingleComment(),
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
	ID          string    // 语言唯一名称，一律小写
	Blocks      []Blocker // 注释块的解析规则定义
	Exts        []string  // 扩展名列表，必须以 . 开头且小写
}

// Get 获取指定语言的定义信息
//
// 若不存在，则返回 nil
func Get(id string) *Language {
	for _, lang := range langs {
		if lang.ID == id {
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
