// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

// 各个语种的语言对照表，通过相应文件的 init() 作初始化这样
// 在删除文件是，就自动删除相应的语言文件，不需要手修改代码。
var locales = map[string]map[string]string{}

// 各个语言需要翻译的所有字符串
const (
	SyntaxError  = "在[%v:%v]出现语法错误[%v]"
	OptionsError = "配置文件[%v]中配置项[%v]错误[%v]"

	// 与 flag 包相关的处理
	FlagUsage = `%v 是一个 RESTful API 文档生成工具。

参数：
%v

源代码采用 MIT 开源许可证，发布于 %v
详细信息可访问官网 %v
`
	FlagHUsage              = "显示帮助信息"
	FlagVUsage              = "显示版本信息"
	FlagLanguagesUsage      = "显示所有支持的语言"
	FlagGUsage              = "创建一个默认的配置文件"
	FlagEncodingsUsage      = "显示支持的编码方式"
	FlagWDUsage             = "指定工作目录，默认为当前目录"
	FlagPprofUsage          = "指定一种调试输出类型，可以为 %s 或是 %s"
	FlagVersionBuildWith    = "%v %v build with %v\n"
	FlagVersionCommitHash   = "commit hash %v\n"
	FlagSupportedLanguages  = "目前支持以下语言 %v\n"
	FlagSupportedEncodings  = "目前支持以下编码 %v\n"
	FlagConfigWritedSuccess = "配置内容成功写入 %v"
	FlagPprofWritedSuccess  = "pprof 的相关数据已经写入到 %v"
	FlagInvalidPprrof       = "无效的 pprof 参数"

	VersionInCompatible = "当前程序与配置文件中指定的版本号不兼容"
	Complete            = "完成！文档保存在：%v，总用时：%v"

	// 错误信息，可能在地方用到
	ErrRequired              = "不能为空"
	ErrInvalidFormat         = "格式不正确"
	ErrDirNotExists          = "目录不存在"
	ErrInvalidBlockType      = "无效的 block.Type 值：%v"
	ErrUnsupportedInputLang  = "无效的输入语言：%v"
	ErrNotFoundEndFlag       = "找不到结束符号"
	ErrNotFoundSupportedLang = "该目录下没有支持的语言文件"
	ErrUnknownTag            = "不认识的标签：%v"
	ErrDuplicateTag          = "重复的标签：%v"
	ErrTagArgTooMuch         = "标签：%v 指定了太多的参数"
	ErrTagArgNotEnough       = "标签：%v 参数不够"
	ErrUnsupportedEncoding   = "不支持的编码方式：%v"
	ErrDirIsEmpty            = "目录下没有需要解析的文件"
	ErrInvalidValue          = "无效的值"

	// logs
	InfoPrefix  = "[INFO] "
	WarnPrefix  = "[WARN] "
	ErrorPrefix = "[ERRO] "
)
