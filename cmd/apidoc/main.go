// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具
//
// 大致的使用方法为：
//  apidoc cmd [args]
// 其中的 cmd 为子命令，args 代码传递给该子命令的参数。
// 可以使用 help 查看每个子命令的具体说明：
//  apidoc help [cmd]
package main

import "github.com/caixw/apidoc/v7/internal/cmd"

func main() {
	cmd.Exec()
}
