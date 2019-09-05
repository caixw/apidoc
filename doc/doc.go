// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

// 表示支持的各种数据类型
const (
	Null    = "null"
	Bool    = "boolean"
	Object  = "object"
	Array   = "array"
	Number  = "number"
	String  = "string"
	Integer = "integer"
)

// Version 表示版本信息，不能独立使用，一般用在其它结构体中
type Version struct {
	Version string `xml:"version,attr,omitempty"`
}

// Reference 引用其它内容
type Reference struct {
	Reference string `xml:"ref,attr,omitempty"`
}

// Deprecated 已经失败的内容，值为失败的版本号
type Deprecated struct {
	Deprecated string `xml:"deprecated,attr,omitempty"`
}

// Doc 文档
type Doc struct {
	XMLName struct{} `xml:"apidoc"`

	APIDoc string `xml:"apidoc,attr"` // 当前的程序版本

	Version string    `xml:"version,attr,omitempty"` // 文档的版本
	Title   string    `xml:"title"`
	Content string    `xml:"content"`
	Contact *Contact  `xml:"contact"`
	License *Link     `xml:"license,omitempty"` // 版本信息
	Tags    []*Tag    `xml:">tag,omitempty"`    // 所有的标签
	Servers []*Server `xml:">server,omitempty"`

	Apis []*API `xml:"-"`
}

// Tag 标签内容
type Tag struct {
	Name        string `xml:"name,attr"`   // 字面名称，需要唯一
	Description string `xml:"description"` // 具体描述
}

// Server 服务信息
type Server struct {
	Name        string `xml:"name,attr"` // 字面名称，需要唯一
	URL         string `xml:"url,attr"`
	Description string `xml:"description,omitempty"` // 具体描述
}

// Contact 描述联系方式
type Contact struct {
	Name  string `xml:"name,attr"`
	URL   string `xml:"url,attr"`
	Email string `xml:"email,attr,omitempty"`
}

// Link 表示一个链接
type Link struct {
	Text string `xml:"text,attr"`
	URL  string `xml:"url,attr"`
}
