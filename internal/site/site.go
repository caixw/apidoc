// SPDX-License-Identifier: MIT

// Package site 将整个 docs 目录当作一个网站进行管理
package site

import "net/http"

// Handler 将 dir 当作一个网站进行管理
//
// dir 应该指向 docs 目录。
func Handler(dir string) http.Handler {
	return http.FileServer(http.Dir(dir))
}
