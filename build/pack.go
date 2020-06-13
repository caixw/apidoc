// SPDX-License-Identifier: MIT

package build

import (
	"encoding/base64"

	"github.com/issue9/errwrap"
	"github.com/issue9/utils"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// 文档内容中可能包含 ` 等特殊字符，如果直接将文档内容以字符串形式保存为 Go 内容，
// 可能造成生成的 Go 文件不是一个合法的 Go 文件。所以采用 base64 对文档进行编码。
var base64Encoding = base64.StdEncoding

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

	var w errwrap.Buffer

	if opt.FileHeader != "" {
		w.Printf("// %s \n\n", opt.FileHeader)
	}

	if opt.Tag != "" {
		w.Printf("// +build %s \n\n", opt.Tag)
	}

	w.Printf("package %s \n\n", opt.PkgName)

	content := base64Encoding.EncodeToString(buf.Bytes())
	w.Printf("const %s = `%s`\n", opt.VarName, content)

	if w.Err != nil {
		return w.Err
	}

	return utils.DumpGoSource(opt.Path, w.Bytes())
}

// Unpack 用于解压由 Pack 输出的内容
func Unpack(buffer string) ([]byte, error) {
	return base64Encoding.DecodeString(buffer)
}

// Pack 将配置文件中指定的文档内容打包成 Go 文件
func (cfg *Config) Pack(opt *PackOptions) error {
	return Pack(cfg.h, opt, cfg.Output, cfg.Inputs...)
}
