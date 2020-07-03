// SPDX-License-Identifier: MIT

package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// 获取 path 的绝对路径
//
// 如果 path 是相对路径的，则将其设置为相对于 wd 的路径。
func abs(path, wd core.URI) (uri core.URI, err error) {
	scheme, p := path.Parse()
	if scheme != "" && scheme != core.SchemeFile {
		return "", locale.NewError(locale.ErrInvalidURIScheme)
	}

	scheme, dir := wd.Parse()
	if scheme != "" && scheme != core.SchemeFile {
		return "", locale.NewError(locale.ErrInvalidURIScheme)
	}
	if !filepath.IsAbs(dir) {
		if dir, err = filepath.Abs(dir); err != nil {
			return "", err
		}
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
		return core.FileURI(filepath.Clean(filepath.Join(dir, ps))), nil
	}
}

// 获取 path 相对于 wd 的路径
//
// 如果两者不存在关联性，则返回 path 的原始值。
// 返回值仅为普通的路径表示，不会带 scheme 内容。
func rel(path, wd core.URI) (uri core.URI, err error) {
	scheme, p := path.Parse()
	if scheme != "" && scheme != core.SchemeFile {
		return "", locale.NewError(locale.ErrInvalidURIScheme)
	}
	if !filepath.IsAbs(p) {
		if p, err = filepath.Abs(p); err != nil {
			return "", err
		}
	}

	scheme, dir := wd.Parse()
	if scheme != "" && scheme != core.SchemeFile {
		return "", locale.NewError(locale.ErrInvalidURIScheme)
	}
	if !filepath.IsAbs(dir) {
		if dir, err = filepath.Abs(dir); err != nil {
			return "", err
		}
	}

	pp, err := filepath.Rel(dir, p)
	if err != nil || strings.HasPrefix(pp, "../") || strings.HasPrefix(pp, "..\\") {
		return path, nil
	}
	return core.URI(pp), nil
}
