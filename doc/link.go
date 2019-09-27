// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Link 表示一个链接
//  <link url="https://example.com">text</link>
type Link struct {
	Text string `xml:",innerxml"`
	URL  string `xml:"url,attr"`
}

// Contact 描述联系方式
//  <contact name="name">
//      <url>https://example.com</url>
//      <email>name@example.com</email>
//  </contact>
type Contact struct {
	Name  string `xml:"name,attr"`
	URL   string `xml:"url,omitempty"`
	Email string `xml:"email,omitempty"`
}

type (
	shadowLink    Link
	shadowContact Contact
)

// UnmarshalXML xml.Unmarshaler
func (l *Link) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	name := "/" + start.Name.Local
	var shadow shadowLink
	if err := d.DecodeElement(&shadow, &start); err != nil {
		return fixedSyntaxError(err, "", name, 0)
	}

	if !is.URL(shadow.URL) {
		return newSyntaxError(name+"#url", locale.ErrInvalidFormat)
	}

	if shadow.Text == "" {
		return newSyntaxError(name, locale.ErrRequired)
	}

	*l = Link(shadow)
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (c *Contact) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	name := "/" + start.Name.Local
	var cc shadowContact
	if err := d.DecodeElement(&cc, &start); err != nil {
		return fixedSyntaxError(err, "", name, 0)
	}

	if cc.Name == "" {
		return newSyntaxError(name+"#name", locale.ErrRequired)
	}

	if cc.URL == "" && cc.Email == "" {
		return newSyntaxError(name+"#url|email", locale.ErrRequired)
	}

	if cc.URL != "" && !is.URL(cc.URL) {
		return newSyntaxError(name+"#url", locale.ErrInvalidFormat)
	}

	if cc.Email != "" && !is.Email(cc.Email) {
		return newSyntaxError(name+"#email", locale.ErrInvalidFormat)
	}

	*c = Contact(cc)
	return nil
}
