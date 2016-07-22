// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

// 之所以用常量代替直接使用字符串，是因为若直接在源码中直接使用字符串，
// 则在修改字符串之后，很难反向查找该字符串所对应的所有语种下的相应代码；
// 而使用常量，则可以通过工具 (godef 等) 可以直接跳转到定义处。

const (
	SyntaxError  = "在[%v:%v]出现语法错误[%v]"      // app/errors.go:23
	OptionsError = "配置文件[%v]中配置项[%v]错误:[%v]" // app/errors.go:27

	FlagHUsage     = "显示帮助信息"                    // main.go:28
	FlagVUsage     = "显示版本信息"                    // main.go:29
	FlagLUsage     = "显示所有支持的语言"                 // main.go:30
	FlagGUsage     = "在当前目录下创建一个默认的配置文件"         // main.go:31
	FlagPprofUsage = "指定一种调试输出类型，可以为 cpu 或是 mem" // main.go:32
)
