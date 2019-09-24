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
	var shadow shadowExample
	if err := d.DecodeElement(&shadow, &start); err != nil {
		return err
	}

	if shadow.Mimetype == "" {
		return locale.Errorf(locale.ErrRequired, "mimetype")
	}

	if shadow.Content == "" {
		return locale.Errorf(locale.ErrRequired, "content")
	}

	*e = Example(shadow)
	return nil
}
