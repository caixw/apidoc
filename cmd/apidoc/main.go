// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具
//
// 大致的使用方法为：
//  apidoc cmd [args]
// 具体的参数说明，可以使用 help 参数查看：
//  apidoc help cmd
package main

import "github.com/caixw/apidoc/v6/internal/cmd"

func main() {
	cmd.Exec()
}
