// SPDX-License-Identifier: MIT

// Package html 处理静态的资源文件
package html

type static struct {
	name        string
	contentType string
	data        []byte
}

// Get 获取指定名称的文件内容以其 content-type
func Get(name string) ([]byte, string) {
	for _, item := range data {
		if item.name == name {
			return item.data, item.contentType
		}
	}

	return nil, ""
}
