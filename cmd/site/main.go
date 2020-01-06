// SPDX-License-Identifier: MIT

// 简单地将 docs 作为一个 web 服务运行
//
// 可作为测试 xsl 使用，访问 localhost:8080/example 测试页面
package main

import (
	"net/http"

	"github.com/caixw/apidoc/v5"
)

func main() {
	http.Handle("/", apidoc.Static("../../docs", false))
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
