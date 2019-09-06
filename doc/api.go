// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

// API 表示 <api> 顶层元素
type API struct {
	XMLName     struct{}    `xml:"api"`
	Version     string      `xml:"version,attr,omitempty"`
	Method      string      `xml:"method,attr"`
	ID          string      `xml:"id,attr,omitempty"`
	Path        *Path       `xml:"path"`
	Description Richtext    `xml:"description,omitempty"`
	Requests    []*Request  `xml:"request"`
	Responses   []*Response `xml:"response"`
	Callback    *Callback   `xml:"callback,omitempty"`

	Tags    []string `xml:"tag,omitempty"`
	Servers []string `xml:"server,omitempty"`
}

// Response 返回的内容
type Response Request

// Request 请求内容
type Request struct {
	Param

	Status   int        `xml:"status,attr"`
	Mimetype string     `xml:"mimetype,attr"`
	Examples []*Example `xml:"example,omitempty"`
	Headers  []*Header  `xml:"header,omitempty"`
}

// Path 路径信息
type Path struct {
	Path      string   `xml:"path,attr"`
	Params    []*Param `xml:">param,omitempty"`
	Queries   []*Param `xml:">query,omitempty"`
	Reference string   `xml:"ref,attr,omitempty"`
}

// Param 表示参数类型
type Param struct {
	Name       string   `xml:"name,attr"`
	Type       string   `xml:"type,attr"`
	Deprecated string   `xml:"deprecated,attr,omitempty"`
	Default    string   `xml:"default,attr,omitempty"`
	Enums      []*Enum  `xml:"enum,omitempty"`
	Array      bool     `xml:"array,attr,omitempty"`
	Items      []*Param `xml:"param,omitempty"`
	Reference  string   `xml:"ref,attr,omitempty"`
}

// Enum 表示枚举值
type Enum struct {
	Deprecated  string   `xml:"deprecated,attr,omitempty"`
	Value       string   `xml:"value,attr"`
	Description Richtext `xml:",innerxml"`
}

// Header 报头信息
type Header struct {
	Name        string   `xml:"name,attr"`
	Description Richtext `xml:",innerxml"`
	Deprecated  string   `xml:"deprecated,attr,omitempty"`
}

// Example 示例代码
type Example struct {
	Type    string   `xml:"type,attr"`
	Content Richtext `xml:",innerxml"`
}

// Callback 回调函数的定义
type Callback struct {
	Response
	Method     string     `xml:"method,attr"`
	Queries    []*Param   `xml:"queries,omitempty"` // 查询参数
	Requests   []*Request `xml:"requests,omitempty"`
	Deprecated string     `xml:"deprecated,attr,omitempty"`
	Reference  string     `xml:"ref,attr,omitempty"`
}
