// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Request 请求内容
type Request struct {
	Type        Type       `xml:"type,attr"`
	Deprecated  Version    `xml:"deprecated,attr,omitempty"`
	Enums       []*Enum    `xml:"enum,omitempty"`
	Array       bool       `xml:"array,attr,omitempty"`
	Items       []*Param   `xml:"param,omitempty"`
	Reference   string     `xml:"ref,attr,omitempty"`
	Summary     string     `xml:"summary,attr,omitempty"`
	Status      Status     `xml:"status,attr,omitempty"`
	Mimetype    string     `xml:"mimetype,attr"`
	Examples    []*Example `xml:"example,omitempty"`
	Headers     []*Header  `xml:"header,omitempty"`
	Description string     `xml:"description,omitempty"`
}

// IsEnum 是否为枚举值
func (r *Request) IsEnum() bool {
	return len(r.Enums) > 0
}

type shadowRequest Request

// UnmarshalXML xml.Unmarshaler
func (r *Request) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var rr shadowRequest
	if err := d.DecodeElement(&rr, &start); err != nil {
		return err
	}

	if rr.Type == None {
		return locale.Errorf(locale.ErrRequired, "type")
	}
	if rr.Type == Object && len(rr.Items) == 0 {
		return locale.Errorf(locale.ErrNeedProperty)
	}

	if rr.Mimetype == "" {
		return locale.Errorf(locale.ErrRequired, "mimetype")
	}

	// 判断 enums 的值是否相同
	if key := getDuplicateEnum(rr.Enums); key != "" {
		return locale.Errorf(locale.ErrDuplicateEnum, key)
	}

	// 判断 items 的值是否相同
	if key := getDuplicateItems(rr.Items); key != "" {
		return locale.Errorf(locale.ErrDuplicateValue, key)
	}

	*r = Request(rr)
	return nil
}
