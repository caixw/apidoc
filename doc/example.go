// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Example 示例代码
type Example struct {
	Description string `xml:"description,omitempty"`
	Mimetype    string `xml:"mimetype,attr"`
	Content     string `xml:",cdata"`
}

type shadowExample Example

// UnmarshalXML xml.Unmarshaler
func (e *Example) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	name := "/" + start.Name.Local
	var shadow shadowExample
	if err := d.DecodeElement(&shadow, &start); err != nil {
		return fixedSyntaxError(err, "", name, 0)
	}

	if shadow.Mimetype == "" {
		return newSyntaxError(name+"#mimetype", locale.ErrRequired)
	}

	if shadow.Content == "" {
		return newSyntaxError(name, locale.ErrRequired)
	}

	*e = Example(shadow)
	return nil
}
