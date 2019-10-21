// SPDX-License-Identifier: MIT

// Package docs docs 内容管理
package docs

import "net/http"

// Handler 将 docs 当作一个网站进行管理
func Handler(docs string) http.Handler {
	return http.FileServer(http.Dir(docs))
}
