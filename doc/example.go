// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Example 示例代码
type Example struct {
	Mimetype string `xml:"mimetype,attr"`
	Content  string `xml:",cdata"`
	Summary  string `xml:"summary,attr,omitempty"`
	Description Richtext `xml:"description,omitempty"`
}

type shadowExample Example

// UnmarshalXML xml.Unmarshaler
func (e *Example) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	shadow := (*shadowExample)(e)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Mimetype == "" {
		return newSyntaxError(field+"#mimetype", locale.ErrRequired)
	}

	if shadow.Content == "" {
		return newSyntaxError(field, locale.ErrRequired)
	}

	return nil
}
