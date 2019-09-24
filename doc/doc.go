// SPDX-License-Identifier: MIT

package doc

import (
	"bytes"
	"encoding/xml"
	"sort"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

// Doc 文档
type Doc struct {
	XMLName struct{} `xml:"apidoc"`

	APIDoc  string    `xml:"apidoc,attr,omitempty"`  // 程序的版本号
	Version Version   `xml:"version,attr,omitempty"` // 文档的版本
	Title   string    `xml:"title"`
	Content string    `xml:"content"`
	Contact *Contact  `xml:"contact"`
	License *Link     `xml:"license,omitempty"` // 版本信息
	Tags    []*Tag    `xml:"tag,omitempty"`     // 所有的标签
	Servers []*Server `xml:"server,omitempty"`
	Apis    []*API    `xml:"api,omitempty"`

	// 应用于全局的变量
	Mimetypes string     `xml:"mimetypes,omitempty"` // 指定可用的 mimetype 类型
	Responses []*Request `xml:"response,omitempty"`
	Requests  []*Request `xml:"Request,omitempty"`

	references map[string]interface{}
	file       string
	line       int
	data       []byte
}

// New 返回 Doc 实例
func New() *Doc {
	return &Doc{
		APIDoc: vars.Version(),
	}
}

type shadowDoc Doc

// UnmarshalXML xml.Unmarshaler
func (doc *Doc) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var obj shadowDoc
	if err := d.DecodeElement(&obj, &start); err != nil {
		line := bytes.Count(doc.data[:d.InputOffset()], []byte{'\n'})
		return message.WithError(doc.file, "", doc.line+line, err)
	}

	apis := doc.Apis
	*doc = Doc(obj)
	doc.Apis = apis
	return nil
}

// FromXML 从 XML 字符串初始化当前的实例
func (doc *Doc) FromXML(file string, line int, data []byte) error {
	doc.file = file
	doc.line = line
	doc.data = data
	return xml.Unmarshal(data, doc)
}

// Sanitize 检测内容是否合法
func (doc *Doc) Sanitize() error {
	doc.APIDoc = vars.Version() // 防止被覆盖

	// Tag.Name 查重
	sort.SliceStable(doc.Tags, func(i, j int) bool {
		return doc.Tags[i].Name > doc.Tags[j].Name
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Tags[i].Name == doc.Tags[i-1].Name {
			return message.NewLocaleError(doc.file, "tag.name", doc.line, locale.ErrDuplicateValue)
		}
	}

	// TODO

	// Server.Name 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].Name > doc.Servers[j].Name
	})
	for i := 1; i < len(doc.Servers); i++ {
		if doc.Servers[i].Name == doc.Servers[i-1].Name {
			return message.NewLocaleError(doc.file, "server.name", doc.line, locale.ErrDuplicateValue)
		}
	}

	// Server.URL 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].URL > doc.Servers[j].URL
	})
	for i := 1; i < len(doc.Servers); i++ {
		if doc.Servers[i].URL == doc.Servers[i-1].URL {
			return message.NewLocaleError(doc.file, "server.url", doc.line, locale.ErrDuplicateValue)
		}
	}

	// 查看 API 中的标签是否都存在
	for _, api := range doc.Apis {
		if err := api.sanitize(); err != nil {
			return err
		}
	}

	return nil
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

func (doc *Doc) requestExists(status Status, mimetype string) bool {
	return doc.requestResponseExists(doc.Requests, status, mimetype)
}

func (doc *Doc) responseExists(status Status, mimetype string) bool {
	return doc.requestResponseExists(doc.Responses, status, mimetype)
}

func (doc *Doc) requestResponseExists(body []*Request, status Status, mimetype string) bool {
	for _, r := range body {
		if r.Status == status && r.Mimetype == mimetype {
			return true
		}
	}

	return false
}
