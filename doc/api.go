// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

// API 表示 <api> 顶层元素
type API struct {
	Version
	Reference

	XMLName     struct{}    `xml:"api"`
	Method      string      `xml:"method,attr"`
	ID          string      `xml:"id,attr,omitempty"`
	Path        *Path       `xml:"path"`
	Description string      `xml:"description,omitempty"`
	Requests    []*Request  `xml:">request"`
	Responses   []*Response `xml:">response"`
	Tags        []*Tag      `xml:">tag,omitempty"`
	Servers     []*Server   `xml:">server,omitempty"`
	Callback    []*Callback `xml:">callback,omitempty"`
}

// Response 返回的内容
type Response Request

// Request 请求内容
type Request struct {
	Param

	Example string    `xml:"example,omitempty"`
	Headers []*Header `xml:">header,omitempty"`
}

// Path 路径信息
type Path struct {
	Path    string   `xml:"path,attr"`
	Params  []*Param `xml:">param,omitempty"`
	Queries []*Param `xml:">query,omitempty"`
}

// Param 表示参数类型
type Param struct {
	Deprecated
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Default string   `xml:"default,attr,omitempty"`
	Enums   []*Enum  `xml:">enum,omitempty"`
	Array   bool     `xml:"array,attr,omitempty"`
	Items   *[]Param `xml:">param,omitempty"`
}

// Enum 表示枚举值
type Enum struct {
	Deprecated
	Value       string `xml:"value,attr"`
	Description string `xml:"description"`
}

// Header 报头信息
type Header struct {
	Deprecated
	Name        string `xml:"name,attr"`
	Description string `xml:"description"`
}

// Callback 回调函数的定义
type Callback struct {
	Deprecated
	Response
	Method   string     `xml:"method,attr"`
	Queries  []*Param   `xml:"queries,omitempty"` // 查询参数
	Requests []*Request `xml:"requests,omitempty"`
}
