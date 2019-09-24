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
	var shadow shadowLink
	if err := d.DecodeElement(&shadow, &start); err != nil {
		return err
	}

	if !is.URL(shadow.URL) {
		return locale.Errorf(locale.ErrInvalidURL, shadow.URL)
	}

	if shadow.Text == "" {
		return locale.Errorf(locale.ErrRequired, "text")
	}

	*l = Link(shadow)
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (c *Contact) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var cc shadowContact
	if err := d.DecodeElement(&cc, &start); err != nil {
		return err
	}

	if cc.Name == "" {
		return locale.Errorf(locale.ErrRequired, "name")
	}

	if cc.URL == "" && cc.Email == "" {
		return locale.Errorf(locale.ErrRequired, "url or email")
	}

	if cc.URL != "" && !is.URL(cc.URL) {
		return locale.Errorf(locale.ErrInvalidURL, cc.URL)
	}

	if cc.Email != "" && !is.Email(cc.Email) {
		return locale.Errorf(locale.ErrInvalidEmail, cc.Email)
	}

	*c = Contact(cc)
	return nil
}
