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
	field := "/" + start.Name.Local
	shadow := (*shadowLink)(l)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if !is.URL(shadow.URL) {
		return newSyntaxError(field+"#url", locale.ErrInvalidFormat)
	}

	if shadow.Text == "" {
		return newSyntaxError(field, locale.ErrRequired)
	}

	return nil
}

// UnmarshalXML xml.Unmarshaler
func (c *Contact) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	shadow := (*shadowContact)(c)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Name == "" {
		return newSyntaxError(field+"#name", locale.ErrRequired)
	}

	if shadow.URL == "" && shadow.Email == "" {
		return newSyntaxError(field+"#url|email", locale.ErrRequired)
	}

	if shadow.URL != "" && !is.URL(shadow.URL) {
		return newSyntaxError(field+"#url", locale.ErrInvalidFormat)
	}

	if shadow.Email != "" && !is.Email(shadow.Email) {
		return newSyntaxError(field+"#email", locale.ErrInvalidFormat)
	}

	return nil
}
