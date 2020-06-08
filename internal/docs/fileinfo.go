// SPDX-License-Identifier: MIT

package docs

import "encoding/base64"

var base64Encoding = base64.StdEncoding

func init() {
	for _, info := range data {
		c, err := base64Encoding.DecodeString(info.Base64)
		if err != nil {
			panic(err)
		}
		info.content = c
	}
}

// FileInfo 被打包文件的信息
type FileInfo struct {
	Name        string // 相对于打包根目录的地址，同时也会被作为路由地址
	ContentType string
	Base64      string
	content     []byte
}

// NewFileInfo 新的 FileInfo 对象
func NewFileInfo(name, ct string, content []byte) *FileInfo {
	return &FileInfo{
		Name:        name,
		ContentType: ct,
		Base64:      base64Encoding.EncodeToString(content),
	}
}
