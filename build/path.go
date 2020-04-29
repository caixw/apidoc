// SPDX-License-Identifier: MIT

package build

import (
	"os"

	"github.com/caixw/apidoc/v7/core"
)

// 获取 path 的绝对路径
//
// 如果 path 是相对路径的，则将其设置为相对于 wd 的路径
func abs(path, wd core.URI) (core.URI, error) {
	if !path.IsNoScheme() { // 包含协议部分，肯定是绝对路径
		return path, nil
	}

	str := string(path)
	switch {
	case str == "":
		return path, nil
	case str[0] == '~':
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		dirURI, err := core.FileURI(dir)
		if err != nil {
			return "", err
		}
		return dirURI.Append(str[1:]), nil
	case str[0] == '.': // 相对路径
		return wd.Append(str[1:]), nil
	case str[0] == '/' || str[0] == os.PathSeparator: // 绝对路径
		return path, nil
	default: // 相对路径
		return wd.Append(str), nil
	}
}
