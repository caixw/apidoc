// SPDX-License-Identifier: MIT

package build

import (
	"github.com/issue9/pack"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// PackOptions 指定了打包文档内容的参数
type PackOptions struct {
	Path       string
	PkgName    string
	VarName    string
	FileHeader string // 在文件头的输出内容，可以为空。
	Tag        string // 指定 +build 指令的标签，可以为空，表示不需要
}

var defaultPackOptions = &PackOptions{
	Path:       "./apidoc/apidoc.go",
	PkgName:    "apidoc",
	VarName:    "APIDOC",
	FileHeader: locale.Sprintf(locale.PackFileHeader, core.Name),
	Tag:        "apidoc",
}

// Pack 将文档内容打包成一个 Go 文件
func Pack(h *core.MessageHandler, opt *PackOptions, o *Output, i ...*Input) error {
	if opt == nil {
		opt = defaultPackOptions
	}

	buf, err := Buffer(h, o, i...)
	if err != nil {
		return err
	}

	return pack.File(buf.String(), opt.PkgName, opt.VarName, opt.FileHeader, opt.Tag, opt.Path)
}

// Unpack 用于解压由 Pack 输出的内容
func Unpack(buffer string) (doc string, err error) {
	err = pack.Unpack(buffer, &doc)
	return
}

// Pack 将配置文件中指定的文档内容打包成 Go 文件
func (cfg *Config) Pack(h *core.MessageHandler, opt *PackOptions) error {
	return Pack(h, opt, cfg.Output, cfg.Inputs...)
}
