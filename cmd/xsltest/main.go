// SPDX-License-Identifier: MIT

// 用于测试 xsl 的展示
//
// 访问 localhost:8080/example/apidoc.xml 测试页面
package main

import (
	"net/http"

	"github.com/caixw/apidoc/v5/internal/site"
)

func main() {
	http.Handle("/", site.Handler("../../docs"))
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
