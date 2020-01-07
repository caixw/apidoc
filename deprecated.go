// SPDX-License-Identifier: MIT

package apidoc

import (
	"time"

	"github.com/caixw/apidoc/v6/input"
	"github.com/caixw/apidoc/v6/message"
	"github.com/caixw/apidoc/v6/output"
)

// Do 解析文档并输出文档内容
//
// Deprecated: 下个版本取消
func (cfg *Config) Do(start time.Time) {
	cfg.Build(start)
}

// Do 解析文档并输出文档内容
//
// Deprecated: 下个版本取消
func Do(h *message.Handler, o *output.Options, i ...*input.Options) error {
	return Build(h, o, i...)
}
