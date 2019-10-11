// SPDX-License-Identifier: MIT

// 用于测试 xsl 的展示
package main

import "net/http"

func main() {
	http.Handle("/", http.FileServer(http.Dir("../../docs")))
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
