// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"sort"
	"strconv"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v6/internal/locale"
)

// Param 表示参数类型
//  <param name="user" deprecated="1.1.1" type="object" array="true">
//      <param name="name" type="string" />
//      <param name="sex" type="string">
//          <enum value="male" summary="male" />
//          <enum value="female" summary="female" />
//      </param>
//      <param name="age" type="number" />
//  </param>
type Param struct {
	XML
	Name        string   `xml:"name,attr"`
	Type        Type     `xml:"type,attr"`
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
	Default     string   `xml:"default,attr,omitempty"`
	Optional    bool     `xml:"optional,attr,omitempty"`
	Array       bool     `xml:"array,attr,omitempty"`
	Items       []*Param `xml:"param,omitempty"`
	Reference   string   `xml:"ref,attr,omitempty"`
	Summary     string   `xml:"summary,attr,omitempty"`
	Enums       []*Enum  `xml:"enum,omitempty"`
	Description Richtext `xml:"description,omitempty"`
}

// Param 转换成 Param 对象
//
// Request 可以说是 Param 的超级，两者在大部分情况下能用。
func (r *Request) Param() *Param {
	if r == nil {
		return nil
	}

	return &Param{
		XML:         r.XML,
		Name:        r.Name,
		Type:        r.Type,
		Deprecated:  r.Deprecated,
		Default:     "",
		Optional:    true,
		Array:       r.Array,
		Items:       r.Items,
		Reference:   r.Reference,
		Summary:     r.Summary,
		Enums:       r.Enums,
		Description: r.Description,
	}
}

// IsEnum 是否为一个枚举类型
func (p *Param) IsEnum() bool {
	return len(p.Enums) > 0
}

type shadowParam Param

// UnmarshalXML xml.Unmarshaler
func (p *Param) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	shadow := (*shadowParam)(p)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Name == "" {
		return newSyntaxError(field+"/@name", locale.ErrRequired)
	}

	if shadow.Type == None {
		return newSyntaxError(field+"/@type", locale.ErrRequired)
	}
	if shadow.Type == Object && len(shadow.Items) == 0 {
		return newSyntaxError(field+"/items", locale.ErrRequired)
	}

	// 判断 enums 的值是否相同
	if key := getDuplicateEnum(shadow.Enums); key != "" {
		return newSyntaxError(field+"/enum", locale.ErrDuplicateValue)
	}

	if err := chkEnumsType(shadow.Type, shadow.Enums, field); err != nil {
		return err
	}

	// 判断 items 的值是否相同
	if key := getDuplicateItems(shadow.Items); key != "" {
		return newSyntaxError(field+"/items", locale.ErrDuplicateValue)
	}

	if err := checkXML(shadow.Array, len(shadow.Items) > 0, &shadow.XML, field); err != nil {
		return err
	}

	if p.Summary == "" && p.Description.Text == "" {
		return newSyntaxError(field+"/summary", locale.ErrRequired)
	}

	return nil
}

// 检测 enums 中的类型是否符合 t 的标准，比如 Number 要求枚举值也都是数值
func chkEnumsType(t Type, enums []*Enum, field string) error {
	if len(enums) == 0 {
		return nil
	}

	switch t {
	case Number:
		for _, enum := range enums {
			if !is.Number(enum.Value) {
				return newSyntaxError(field+"/enum/@"+enum.Value, locale.ErrInvalidFormat)
			}
		}
	case Bool:
		for _, enum := range enums {
			if _, err := strconv.ParseBool(enum.Value); err != nil {
				return newSyntaxError(field+"/enum/@"+enum.Value, locale.ErrInvalidFormat)
			}
		}
	case Object, None:
		return newSyntaxError(field+"/enum", locale.ErrInvalidValue)
	}

	return nil
}

// 返回重复枚举的值
func getDuplicateEnum(enums []*Enum) string {
	if len(enums) == 0 {
		return ""
	}

	es := make([]string, 0, len(enums))
	for _, e := range enums {
		es = append(es, e.Value)
	}
	sort.SliceStable(es, func(i, j int) bool { return es[i] > es[j] })

	for i := 1; i < len(es); i++ {
		if es[i] == es[i-1] {
			return es[i]
		}
	}

	return ""
}

func getDuplicateItems(items []*Param) string {
	if len(items) == 0 {
		return ""
	}

	es := make([]string, 0, len(items))
	for _, e := range items {
		es = append(es, e.Name)
	}
	sort.SliceStable(es, func(i, j int) bool { return es[i] > es[j] })

	for i := 1; i < len(es); i++ {
		if es[i] == es[i-1] {
			return es[i]
		}
	}

	return ""
}
