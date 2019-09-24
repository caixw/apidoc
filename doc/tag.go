// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Tag 标签内容
//  <tag name="tag1" deprecated="1.1.1">description</tag>
type Tag struct {
	Name        string  `xml:"name,attr"` // 字面名称，需要唯一
	Deprecated  Version `xml:"deprecated,attr,omitempty"`
	Description string  `xml:",cdata"`
}

// Server 服务信息
//  <server name="tag1" deprecated="1.1.1" url="api.example.com/admin">description</server>
type Server struct {
	Name        string  `xml:"name,attr"` // 字面名称，需要唯一
	URL         string  `xml:"url,attr"`
	Deprecated  Version `xml:"deprecated,attr,omitempty"`
	Description string  `xml:",cdata"`
}

type (
	shadowTag    Tag
	shadowServer Server
)

// UnmarshalXML xml.Unmarshaler
func (t *Tag) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var tt shadowTag
	if err := d.DecodeElement(&tt, &start); err != nil {
		return err
	}

	if tt.Name == "" {
		return locale.Errorf(locale.ErrRequired, "name")
	}

	if tt.Description == "" {
		return locale.Errorf(locale.ErrRequired, "description")
	}

	*t = Tag(tt)
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (srv *Server) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var ss shadowServer
	if err := d.DecodeElement(&ss, &start); err != nil {
		return err
	}

	if ss.Name == "" {
		return locale.Errorf(locale.ErrRequired, "name")
	}

	if ss.Description == "" {
		return locale.Errorf(locale.ErrRequired, "description")
	}

	if !is.URL(ss.URL) {
		return locale.Errorf(locale.ErrInvalidURL, ss.URL)
	}

	*srv = Server(ss)
	return nil
}
