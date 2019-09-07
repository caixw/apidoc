// SPDX-License-Identifier: MIT

package doc

// Richtext 富文本内容
type Richtext string

// Doc 文档
type Doc struct {
	XMLName struct{} `xml:"apidoc"`

	APIDoc string `xml:"-"` // 程序的版本号

	Version Version   `xml:"version,attr,omitempty"` // 文档的版本
	Title   string    `xml:"title"`
	Content Richtext  `xml:"content"`
	Contact *Contact  `xml:"contact"`
	License *Link     `xml:"license,omitempty"` // 版本信息
	Tags    []*Tag    `xml:"tag,omitempty"`     // 所有的标签
	Servers []*Server `xml:"server,omitempty"`
	Apis    []*API    `xml:"apis,omitempty"`

	// TODO 应用于全局的变量
	//Responses []*Response `xml:"response,omitempty"`
	//Requests  []*Request  `xml:"Request,omitempty"`
	//Mimetypes string `` // 指定可用的 mimetype 类型

	references map[string]interface{}
	file       string
	line       int
}

// Tag 标签内容
type Tag struct {
	Name        string   `xml:"name,attr"`  // 字面名称，需要唯一
	Description Richtext `xml:",omitempty"` // 具体描述
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
}

// Server 服务信息
type Server struct {
	Name        string   `xml:"name,attr"` // 字面名称，需要唯一
	URL         string   `xml:"url,attr"`
	Description Richtext `xml:",omitempty"` // 具体描述
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
}

// Contact 描述联系方式
type Contact struct {
	Name  string `xml:"name,attr"`
	URL   string `xml:"url"`
	Email string `xml:"email,omitempty"`
}

// Link 表示一个链接
type Link struct {
	Text string `xml:",innerxml"`
	URL  string `xml:"url,attr"`
}
