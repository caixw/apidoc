// SPDX-License-Identifier: MIT

package spec

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Example 示例代码
type Example struct {
	Mimetype    string   `xml:"mimetype,attr"`
	Content     string   `xml:",cdata"`
	Summary     string   `xml:"summary,attr,omitempty"`
	Description Richtext `xml:"description,omitempty"`
}

type shadowExample Example

// UnmarshalXML xml.Unmarshaler
func (e *Example) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	shadow := (*shadowExample)(e)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(core.Location{}, err, field)
	}

	if shadow.Mimetype == "" {
		return newSyntaxError(core.Location{}, field+"/@mimetype", locale.ErrRequired)
	}

	if shadow.Content == "" {
		return newSyntaxError(core.Location{}, field, locale.ErrRequired)
	}

	return nil
}
