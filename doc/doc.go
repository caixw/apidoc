// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/vars"
)

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

	// 应用于全局的变量
	Mimetypes string     `xml:"mimetypes,omitempty"` // 指定可用的 mimetype 类型
	Responses []*Request `xml:"response,omitempty"`
	Requests  []*Request `xml:"Request,omitempty"`

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
	Tag
	URL string `xml:"url,attr"`
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

// New 返回 Doc 实例
func New() *Doc {
	return &Doc{
		APIDoc: vars.Version(),
	}
}

// FromXML 从 XML 字符串初始化当前的实例
func (doc *Doc) FromXML(data []byte) error {
	return xml.Unmarshal(data, doc)
}

func (doc *Doc) tagExists(tag string) bool {
	for _, s := range doc.Tags {
		if s.Name == tag {
			return true
		}
	}
	return false
}

func (doc *Doc) serverExists(srv string) bool {
	for _, s := range doc.Servers {
		if s.Name == srv {
			return true
		}
	}
	return false
}

func (doc *Doc) requestExists(status int, mimetype string) bool {
	return doc.requestResponseExists(doc.Requests, status, mimetype)
}

func (doc *Doc) responseExists(status int, mimetype string) bool {
	return doc.requestResponseExists(doc.Responses, status, mimetype)
}

func (doc *Doc) requestResponseExists(body []*Request, status int, mimetype string) bool {
	for _, r := range body {
		if r.Status == status && r.Mimetype == mimetype {
			return true
		}
	}

	return false
}
