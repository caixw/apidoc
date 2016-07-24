// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"io"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// 保证有个初始化的值，部分包的测试功能依赖此变量
var localePrinter *message.Printer = message.NewPrinter(language.Chinese)

// 各个语种的语言对照表，通过相应文件的 init() 作初始化这样
// 在删除文件是，就自动删除相应的语言文件，不需要手修改代码。
var locales = map[string]map[string]string{}

// 各个语言需要翻译的所有字符串
const (
	SyntaxError  = "在[%v:%v]出现语法错误[%v]"     // app/errors.go:23
	OptionsError = "配置文件[%v]中配置项[%v]错误[%v]" // app/errors.go:27

	// 与 flag 包相关的处理
	FlagUsage = `%v 是一个 RESTful API 文档生成工具。

参数：
%v

源代码采用 MIT 开源许可证，发布于 %v
详细信息可访问官网 %v
`
	FlagHUsage              = "显示帮助信息"                    // main.go:28
	FlagVUsage              = "显示版本信息"                    // main.go:29
	FlagLUsage              = "显示所有支持的语言"                 // main.go:30
	FlagGUsage              = "在当前目录下创建一个默认的配置文件"         // main.go:31
	FlagPprofUsage          = "指定一种调试输出类型，可以为 cpu 或是 mem" // main.go:32
	FlagVersionBuildWith    = "%v %v build with %v\n"     // main.go:41
	FlagSupportedLangs      = "目前支持以下语言 %v\n"             // main.go:44
	FlagConfigWritedSuccess = "配置内容成功写入 %v"               // main.go:56
	FlagPprofWritedSuccess  = "pprof 的相关数据已经写入到 %v"       // main.go:73
	FlagInvalidPprrof       = "无效的 pprof 参数"              // main.go:89

	VersionInCompatible = "当前程序与配置文件中指定的版本号不兼容" // main.go:131
	Complete            = "完成！文档保存在：%v，总用时：%v"  // main.go:160

	DebugPort     = "当前为模板调试模式，调试端口为：%v" // output/html.go:58
	DebugTemplate = "当前为模板调试模式，调试模板为：%v" // output/html.go:59

	// 错误信息，可能在地方用到
	ErrRequired              = "不能为空"
	ErrInvalidFormat         = "格式不正确"
	ErrDirNotExists          = "目录不存在"
	ErrInvalidOutputType     = "无效的输出类型"                  // output/output.go
	ErrTemplateNotExists     = "模板不存在"                    // output/output.go
	ErrMkdirError            = "创建目录时发生以下错误：%v"           // output/output.go:51
	ErrInvalidBlockType      = "无效的 block.Type 值：%v"      // input/block
	ErrUnsupportedInputLang  = "无效的输入语言：%v"               // input
	ErrNotFoundEndFlag       = "找不到结束符号"                  // input
	ErrNotFoundSupportedLang = "该目录下没有支持的语言文件"            // input/lang.go
	ErrUnknownTopTag         = "不认识的顶层标签：%v"              // doc/api.go:28
	ErrUnknownTag            = "不认识的标签：%v"                // doc/api.go
	ErrDuplicateTag          = "重复的标签：%v"                 // doc
	ErrTagArgTooMuch         = "标签：%v 指定了太多的参数"           // doc
	ErrTagArgNotEnough       = "标签：%v 参数不够"               // doc
	ErrSecondArgMustURL      = "@apiLicense 第二个参数必须为 URL" // doc
)

// Printer 获取当前语言的 *message.Printer 实例
func Printer() *message.Printer {
	return localePrinter
}

// Print 类型 fmt.Print，与特定的语言绑定。
func Print(v ...interface{}) (int, error) {
	return localePrinter.Print(v...)
}

// Println 类型 fmt.Println，与特定的语言绑定。
func Println(v ...interface{}) (int, error) {
	return localePrinter.Println(v...)
}

// Printf 类型 fmt.Printf，与特定的语言绑定。
func Printf(key string, v ...interface{}) (int, error) {
	return localePrinter.Printf(key, v...)
}

// Sprint 类型 fmt.Sprint，与特定的语言绑定。
func Sprint(v ...interface{}) string {
	return localePrinter.Sprint(v...)
}

// Sprintln 类型 fmt.Sprintln，与特定的语言绑定。
func Sprintln(v ...interface{}) string {
	return localePrinter.Sprintln(v...)
}

// Sprintf 类型 fmt.Sprintf，与特定的语言绑定。
func Sprintf(key message.Reference, v ...interface{}) string {
	return localePrinter.Sprintf(key, v...)
}

// Fprint 类型 fmt.Fprint，与特定的语言绑定。
func Fprint(w io.Writer, v ...interface{}) (int, error) {
	return localePrinter.Fprint(w, v...)
}

// Fprintln 类型 fmt.Fprintln，与特定的语言绑定。
func Fprintln(w io.Writer, v ...interface{}) (int, error) {
	return localePrinter.Fprintln(w, v...)
}

// Fprintf 类型 fmt.Fprintf，与特定的语言绑定。
func Fprintf(w io.Writer, key message.Reference, v ...interface{}) (int, error) {
	return localePrinter.Fprintf(w, key, v...)
}
