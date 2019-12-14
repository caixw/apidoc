// SPDX-License-Identifier: MIT

package apidoc

import (
	"bytes"
	"time"

	"github.com/caixw/apidoc/v5/message"
)

// Make 根据 wd 目录下的配置文件生成文档
//
// Deprecated: 下个版本将弃用，请使用 Config.Do 方法
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
// Deprecated: 下个版本将弃用，请使用 Config.Buffer 方法
func MakeBuffer(h *message.Handler, wd string) *bytes.Buffer {
	cfg := LoadConfig(h, wd)
	return cfg.Buffer()
}

