// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具
//
// 大致的使用方法为：
//  apidoc [options] [path]
// 具体的参数说明，可以使用 h 参数查看：
//  apidoc -h
// path 表示目录列表，多个目录使用空格分隔。
// 用于在 path 下生成配置文件或是从 path 目录加载配置文件。
package main

import "github.com/caixw/apidoc/v5/internal/cmd"

func main() {
	cmd.Exec()
}
