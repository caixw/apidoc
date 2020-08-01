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
		blocks:      cStyle,
	},

	{
		DisplayName: "C/C++",
		ID:          "c++",
		Exts:        []string{".h", ".c", ".cpp", ".cxx", ".hpp"},
		blocks:      cStyle,
	},

	{
		DisplayName: "D",
		ID:          "d",
		Exts:        []string{".d"},
		blocks:      cStyle,
	},

	{
		DisplayName: "Dart",
		ID:          "dart",
		Exts:        []string{".dart"},
		blocks: []blocker{
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
		blocks: []blocker{
			newCStyleString(),
			newSingleComment("%"),
		},
	},

	{
		DisplayName: "Go",
		ID:          "go",
		Exts:        []string{".go"},
		blocks: []blocker{
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
		blocks: []blocker{
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
		blocks:      cStyle,
	},

	{
		DisplayName: "JavaScript",
		ID:          "javascript",
		Exts:        []string{".js"},
		blocks: []blocker{
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
		DisplayName: "Julia",
		ID:          "julia",
		Exts:        []string{".jl"},
		blocks: []blocker{
			newString("'", "'", `\`),
			newString(`r"`, `"`, ""),
			newString(`b"`, `"`, ""),
			newString(`v"`, `"`, ""),
			newString(`raw"`, `"`, ""),
			newString(`"""`, `"""`, ``),
			newString(`"`, `"`, `\`),
			newMultipleComment("#=", "=#", "#="),
			newSingleComment("#"),
		},
	},

	{
		DisplayName: "Kotlin",
		ID:          "kotlin",
		Exts:        []string{".kt"},
		blocks:      cStyle,
	},

	{
		DisplayName: "Lisp/Clojure",
		ID:          "lisp",
		// .ss,.scm ==> scheme
		// .clj     ==> Clojure
		Exts: []string{".lisp", ".lsp", ".l", ".ss", ".scm", ".clj"},
		blocks: []blocker{
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
		blocks: []blocker{
			newString("'", "'", `\`),
			newString("\"", "\"", `\`),
			newString("[[", "]]", ``),
			newMultipleComment("--[====[", "]====]", "-="),
			newMultipleComment("--[[", "]]", "-="),
			newSingleComment("--"), // 放在 --[[ 之后，否则会把 --[[ 当作 -- 解析
		},
	},

	{
		DisplayName: "Nim",
		ID:          "nim",
		Exts:        []string{".nim"},
		blocks: []blocker{
			newString("'", "'", "\\"),
			newString(`"`, `"`, "\\"),
			newNimRawString(),
			newNimMultipleString(),
			newSwiftNestMCommentBlock("##[", "]##", "#"),
			newSwiftNestMCommentBlock("#[", "]#", "#"),
			newSingleComment("#"),
		},
	},

	{
		DisplayName: "Pascal/Delphi",
		ID:          "pascal",
		Exts:        []string{".pas", ".pp"},
		blocks: []blocker{
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
		blocks: []blocker{
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
		blocks: []blocker{
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
		blocks: []blocker{
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
		blocks: []blocker{
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
		blocks:      cStyle,
	},

	{
		DisplayName: "Scala",
		ID:          "scala",
		Exts:        []string{".scala"},
		blocks:      cStyle,
	},

	{
		DisplayName: "Swift",
		ID:          "swift",
		Exts:        []string{".swift"},
		blocks: []blocker{
			newString(`"""`, `"""`, `\`),
			newString(`#"`, `"#`, ""),
			newCStyleString(),
			newCStyleChar(),
			newCStyleSingleComment(),
			newSwiftNestMCommentBlock("/*", "*/", "*"),
		},
	},

	{
		DisplayName: "TypeScript",
		ID:          "typescript",
		Exts:        []string{".ts"},
		blocks: []blocker{
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
		blocks: []blocker{
			newCStyleString(),
			newCStyleChar(),
			newSingleComment("///"), // 需要在 // 之前定义
			newSingleComment("//!"), // 需要在 // 之前定义
			newCStyleSingleComment(),
		},
	},
}

var cStyle = []blocker{
	newCStyleString(),
	newCStyleChar(),
	newSingleComment("///"), // 需要在 // 之前定义
	newCStyleSingleComment(),
	newCStyleMultipleComment(),
}

// 处理 "XXX\""
func newCStyleString() blocker {
	return newString(`"`, `"`, `\`)
}

// 处理 '"'
func newCStyleChar() blocker {
	return newString(`'`, `'`, "")
}

// 处理 // xxx
func newCStyleSingleComment() blocker {
	return newSingleComment(`//`)
}

// 处理 /* */
func newCStyleMultipleComment() blocker {
	return newMultipleComment(`/*`, `*/`, "*")
}

// Language 语言模块的定义
type Language struct {
	DisplayName string    // 显示友好的名称
	ID          string    // 语言唯一名称，一律小写
	blocks      []blocker // 注释块的解析规则定义
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
