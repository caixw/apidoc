// SPDX-License-Identifier: MIT

package build

import (
	"os"
	"path/filepath"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// 获取 path 的绝对路径
//
// 如果 path 是相对路径的，则将其设置为相对于 wd 的路径
func abs(path, wd core.URI) (core.URI, error) {
	scheme, p := path.Parse()
	if scheme != "" && scheme != core.SchemeFile {
		return "", locale.NewError(locale.ErrInvalidURIScheme)
	}

	scheme, dir := wd.Parse()
	if scheme != "" && scheme != core.SchemeFile {
		return "", locale.NewError(locale.ErrInvalidURIScheme)
	}

	ps := string(p)
	switch {
	case ps == "":
		return core.FileURI(dir), nil
	case ps[0] == '~':
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return core.FileURI(dir).Append(ps[1:]), nil
	case filepath.IsAbs(ps) || ps[0] == '/' || ps[0] == os.PathSeparator:
		return core.FileURI(ps), nil
	default: // 相对路径
		p := filepath.Clean(filepath.Join(dir, ps))
		return core.FileURI(p), nil
	}
}
