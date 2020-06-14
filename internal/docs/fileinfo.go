// SPDX-License-Identifier: MIT

package docs

import "github.com/issue9/pack"

// FileInfo 被打包文件的信息
type FileInfo struct {
	Name        string // 相对于打包根目录的地址，同时也会被作为路由地址
	ContentType string
	Content     []byte
}

func files() []*FileInfo {
	var fis []*FileInfo
	if err := pack.Unpack(data, &fis); err != nil {
		panic(err)
	}

	return fis
}
