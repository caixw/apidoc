// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Header 报头信息
type Header struct {
	Name        string  `xml:"name,attr"`
	Deprecated  Version `xml:"deprecated,attr,omitempty"`
	Description string  `xml:",cdata"`
}

type shadowHeader Header

// UnmarshalXML xml.Unmarshaler
func (h *Header) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	shadow := (*shadowHeader)(h)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Name == "" {
		return newSyntaxError(field+"#name", locale.ErrRequired)
	}

	if shadow.Description == "" {
		return newSyntaxError(field, locale.ErrRequired)
	}

	return nil
}
