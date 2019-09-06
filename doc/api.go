// SPDX-License-Identifier: MIT

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
	Deprecated  string      `xml:"deprecated,attr,omitempty"`

	Tags    []string `xml:"tag,omitempty"`
	Servers []string `xml:"server,omitempty"`

	line int
	file string
}

// Response 返回的内容
type Response Request

// Request 请求内容
type Request struct {
	Param

	Status      int        `xml:"status,attr"`
	Mimetype    string     `xml:"mimetype,attr"`
	Examples    []*Example `xml:"example,omitempty"`
	Headers     []*Header  `xml:"header,omitempty"`
	Description Richtext   `xml:",innerxml"`
}

// Path 路径信息
type Path struct {
	Path      string   `xml:"path,attr"`
	Params    []*Param `xml:"param,omitempty"`
	Queries   []*Param `xml:"query,omitempty"`
	Reference string   `xml:"ref,attr,omitempty"`
}

// Param 表示参数类型
type Param struct {
	Name       string   `xml:"name,attr"`
	Type       string   `xml:"type,attr"`
	Deprecated string   `xml:"deprecated,attr,omitempty"`
	Default    string   `xml:"default,attr,omitempty"`
	Required   bool     `xml:"required,attr,omitempty"`
	Enums      []*Enum  `xml:"enum,omitempty"`
	Array      bool     `xml:"array,attr,omitempty"`
	Items      []*Param `xml:"param,omitempty"`
	Reference  string   `xml:"ref,attr,omitempty"`
	Summary    string   `xml:"summary,attr,omitempty"`
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
	Mimetype string   `xml:"mimetype,attr"`
	Content  Richtext `xml:",innerxml"`
}

// Callback 回调函数的定义
type Callback struct {
	Param
	Schema      string     `xml:"schema,attr"` // http 或是 https
	Description Richtext   `xml:",innerxml"`
	Mimetype    string     `xml:"mimetype,attr"`
	Examples    []*Example `xml:"example,omitempty"`
	Headers     []*Header  `xml:"header,omitempty"`
	Method      string     `xml:"method,attr"`
	Queries     []*Param   `xml:"queries,omitempty"` // 查询参数
	Requests    []*Request `xml:"requests,omitempty"`
	Deprecated  string     `xml:"deprecated,attr,omitempty"`
	Reference   string     `xml:"ref,attr,omitempty"`

	// 对回调的返回要求
	Responses []*Response `xml:"response,omitempty"`
}
