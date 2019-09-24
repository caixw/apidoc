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
	var hh shadowHeader
	if err := d.DecodeElement(&hh, &start); err != nil {
		return err
	}

	if hh.Name == "" {
		return locale.Errorf(locale.ErrRequired, "name")
	}

	if hh.Description == "" {
		return locale.Errorf(locale.ErrRequired, "description")
	}

	*h = Header(hh)
	return nil
}
