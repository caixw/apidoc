// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Request 请求内容
type Request struct {
	XML

	// 一般无用，但是用于描述 XML 对象时，可以用来表示顶层元素的名称
	Name string `xml:"name,attr,omitempty"`

	Type        Type       `xml:"type,attr"`
	Deprecated  Version    `xml:"deprecated,attr,omitempty"`
	Enums       []*Enum    `xml:"enum,omitempty"`
	Array       bool       `xml:"array,attr,omitempty"`
	Items       []*Param   `xml:"param,omitempty"`
	Reference   string     `xml:"ref,attr,omitempty"`
	Summary     string     `xml:"summary,attr,omitempty"`
	Status      Status     `xml:"status,attr,omitempty"`
	Mimetype    string     `xml:"mimetype,attr,omitempty"`
	Examples    []*Example `xml:"example,omitempty"`
	Headers     []*Param   `xml:"header,omitempty"` // 当前独有的报头，公用的可以放在 API 中
	Description Richtext   `xml:"description,omitempty"`
}

// IsEnum 是否为枚举值
func (r *Request) IsEnum() bool {
	return len(r.Enums) > 0
}

type shadowRequest Request

// UnmarshalXML xml.Unmarshaler
func (r *Request) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	shadow := (*shadowRequest)(r)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Type == Object && len(shadow.Items) == 0 {
		return newSyntaxError(field+"/param", locale.ErrRequired)
	}

	// 判断 enums 的值是否相同
	if key := getDuplicateEnum(shadow.Enums); key != "" {
		return newSyntaxError(field+"/enum", locale.ErrDuplicateValue)
	}

	if err := chkEnumsType(shadow.Type, shadow.Enums, field); err != nil {
		return err
	}

	if err := checkXML(&shadow.XML, field); err != nil {
		return err
	}

	// 判断 items 的值是否相同
	if key := getDuplicateItems(shadow.Items); key != "" {
		return newSyntaxError(field+"/param", locale.ErrDuplicateValue)
	}

	return nil
}
