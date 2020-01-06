// SPDX-License-Identifier: MIT

package apidoc

import (
	"bytes"
	"net/http"
	"time"

	"github.com/caixw/apidoc/v5/message"
)

// Make 根据 wd 目录下的配置文件生成文档
//
// Deprecated: 下个版本将弃用，请使用 Config.Do 代替。
func Make(h *message.Handler, wd string, test bool) {
	now := time.Now()

	cfg := LoadConfig(h, wd)
	if test {
		cfg.Test()
		return
	}
	cfg.Do(now)
}

// MakeBuffer 根据 wd 目录下的配置文件生成文档内容并保存至内存
//
// Deprecated: 下个版本将弃用，请使用 Config.Buffer 代替。
func MakeBuffer(h *message.Handler, wd string) *bytes.Buffer {
	cfg := LoadConfig(h, wd)
	return cfg.Buffer()
}

// Site 返回文件服务中间件
//
// Deprecated: 下个版本弃用，请使用 Static 代替。
func Site(dir string) http.Handler {
	return Static(dir, false)
}
