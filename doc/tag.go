// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"sort"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Tag 标签内容
//  <tag name="tag1" deprecated="1.1.1" />
type Tag struct {
	Name       string  `xml:"name,attr"`  // 标签的唯一 ID
	Title      string  `xml:"title,attr"` // 显示的名称
	Deprecated Version `xml:"deprecated,attr,omitempty"`
}

// Server 服务信息
//  <server name="tag1" deprecated="1.1.1" url="api.example.com/admin" summary="description" />
type Server struct {
	Name        string   `xml:"name,attr"` // 字面名称，需要唯一
	URL         string   `xml:"url,attr"`
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
	Summary     string   `xml:"summary,attr,omitempty"`
	Description Richtext `xml:"description,omitempty"`
}

type (
	shadowTag    Tag
	shadowServer Server
)

// UnmarshalXML xml.Unmarshaler
func (t *Tag) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	var shadow shadowTag
	if err := d.DecodeElement(&shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Name == "" {
		return newSyntaxError(field+"/@name", locale.ErrRequired)
	}

	if shadow.Title == "" {
		return newSyntaxError(field+"/@title", locale.ErrRequired)
	}

	*t = Tag(shadow)
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (srv *Server) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	var shadow shadowServer
	if err := d.DecodeElement(&shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Name == "" {
		return newSyntaxError(field+"/@name", locale.ErrRequired)
	}

	if !is.URL(shadow.URL) {
		return newSyntaxError(field+"/@url", locale.ErrInvalidFormat)
	}

	if shadow.Summary == "" && shadow.Description.Text == "" {
		return newSyntaxError(field+"/summary", locale.ErrRequired)
	}

	*srv = Server(shadow)
	return nil
}

// 查找是否有重复的标签
//
// 根据标签名称查找
func findDupTag(tags []*Tag) string {
	ts := make([]string, 0, len(tags))
	for _, tag := range tags {
		ts = append(ts, tag.Name)
	}
	return findDupString(ts)
}

// 查找是否有重复的 Server 标签
//
// 根据 name 和 url 查找
func findDupServer(servers []*Server) string {
	names := make([]string, 0, len(servers))
	urls := make([]string, 0, len(servers))
	for _, srv := range servers {
		names = append(names, srv.Name)
		urls = append(urls, srv.URL)
	}

	if key := findDupString(names); key != "" {
		return key
	}
	return findDupString(urls)
}

func findDupString(keys []string) string {
	sort.Strings(keys)

	for i := 1; i < len(keys); i++ {
		if keys[i-1] == keys[i] {
			return keys[i]
		}
	}

	return ""
}
