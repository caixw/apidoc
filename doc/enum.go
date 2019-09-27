// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Enum 表示枚举值
//  <enum value="male">男性</enum>
//  <enum value="female">女性</enum>
type Enum struct {
	Deprecated  Version `xml:"deprecated,attr,omitempty"`
	Value       string  `xml:"value,attr"`
	Description string  `xml:",cdata"`
}

type shadowEnum Enum

// UnmarshalXML xml.Unmarshaler
func (e *Enum) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	name := "/" + start.Name.Local
	var shadow shadowEnum
	if err := d.DecodeElement(&shadow, &start); err != nil {
		return fixedSyntaxError(err, "", name, 0)
	}

	if shadow.Value == "" {
		return newSyntaxError(name+"#value", locale.ErrRequired)
	}

	if shadow.Description == "" {
		return newSyntaxError(name, locale.ErrRequired)
	}

	*e = Enum(shadow)
	return nil
}
