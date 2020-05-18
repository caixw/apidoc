// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具
//
// 大致的使用方法为：
//  apidoc cmd [args]
// 其中的 cmd 为子命令，args 代码传递给该子命令的参数。
// 可以使用 help 查看每个子命令的具体说明：
//  apidoc help [cmd]
package main

import (
	"fmt"
	"os"

	"github.com/issue9/utils"
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/internal/cmd"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func main() {
	tag, err := utils.GetSystemLanguageTag()
	if err != nil { // 无法获取系统语言，则采用默认值
		fmt.Fprintln(os.Stderr, err, tag)
		tag = language.MustParse(locale.DefaultLocaleID)
	}

	if err := cmd.Init(os.Stdout, tag).Exec(os.Args[1:]); err != nil {
		panic(err)
	}
}
