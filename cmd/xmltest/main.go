// SPDX-License-Identifier: MIT

// 用于 xml 的测试命令
package main

import (
	"net/http"

	"github.com/caixw/apidoc/v5"
)

func main() {
	http.Handle("/apidoc/", apidoc.Handle("../../docs/apidoc.xml", "", nil))
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
