// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func validXML(p *ast.Request, content []byte) error {
	if len(content) == 0 {
		if p == nil || p.Type.V() == ast.TypeNone {
			return nil
		}
		return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
	}

	d := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := d.Token()
		if err == io.EOF && token == nil { // 正常结束
			return nil
		}
		if err != nil {
			return err
		}

		switch v := token.(type) {
		case xml.StartElement:
			return validElement(v, p.Param(), true, d, v.Name.Local)
		case xml.EndElement:
			return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
		}
	}
}

func validElement(start xml.StartElement, p *ast.Param, allowArray bool, d *xml.Decoder, field string) error {
	if err := validStartElement(start, p, allowArray, d, field); err != nil {
		return err
	}

	for {
		token, err := d.Token()
		if err == io.EOF && token == nil { // 正常结束
			return nil
		}
		if err != nil {
			return err
		}

		var chardata []byte
		switch v := token.(type) {
		case xml.StartElement:
			if allowArray && p.Array.V() && v.Name.Local == p.Name.V() {
				return validElement(v, p, false, d, buildXMLField(field, p))
			}

			for _, pp := range p.Items {
				if pp.Name.V() == v.Name.Local {
					if pp.XMLExtract.V() {
						return validElement(v, p, true, d, buildXMLField(field, pp))
					}
					return validElement(v, pp, true, d, buildXMLField(field, p))
				}
			}
		case xml.EndElement:
			if allowArray && p.Array.V() {
				if p.XMLWrapped.V() != "" && p.XMLWrapped.V() != v.Name.Local {
					return core.NewSyntaxError(core.Location{}, v.Name.Local, locale.ErrNotFound)
				}
				return nil
			}
			if p.Name.V() != v.Name.Local {
				return core.NewSyntaxError(core.Location{}, v.Name.Local, locale.ErrNotFound)
			}

			if chardata != nil {
				return validXMLParamValue(p, p.Name.V(), string(chardata))
			}
			return nil
		case xml.CharData:
			chardata = v
		}
	}
}

func validStartElement(start xml.StartElement, p *ast.Param, allowArray bool, d *xml.Decoder, field string) error {
	for _, attr := range start.Attr {
		for _, pp := range p.Items {
			if attr.Name.Local != pp.Name.V() {
				continue
			}
			if err := validXMLParamValue(pp, buildXMLField(field, pp), attr.Value); err != nil {
				return err
			}
			break
		}
	}

	if allowArray && p.Array.V() {
		if p.XMLWrapped.V() != "" && p.XMLWrapped.V() != start.Name.Local {
			return core.NewSyntaxError(core.Location{}, start.Name.Local, locale.ErrNotFound)
		}

		start.Attr = start.Attr[:0] // 在 allowArray == true 已经处理过 start.Attr
		return validElement(start, p, false, d, buildXMLField(field, p))
	}

	if (p.Name.V() != start.Name.Local) && (!allowArray && p.XMLWrapped.V() != start.Name.Local) {
		return core.NewSyntaxError(core.Location{}, start.Name.Local, locale.ErrNotFound)
	}
	return nil
}

func buildXMLField(field string, p *ast.Param) string {
	if p.XMLAttr.V() {
		return field + "@" + p.Name.V()
	}
	return field + p.Name.V()
}

// 验证 p 描述的类型与 v 是否匹配，如果不匹配返回错误信息。
// field 表示 p 在整个对象中的位置信息。
func validXMLParamValue(p *ast.Param, field, v string) error {
	switch p.Type.V() {
	case ast.TypeNone:
		if v != "" {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidValue)
		}
	case ast.TypeNumber:
		if !is.Number(v) {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeBool:
		if _, err := strconv.ParseBool(v); err != nil {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeString, ast.TypeObject:
		return nil
	default:
		panic(fmt.Sprintf("文档中类型定义错误 %s", p.Type.V()))
	}

	if isEnum(p) {
		for _, enum := range p.Enums {
			if enum.Value.V() == v {
				return nil
			}
		}
		return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidValue)
	}

	return nil
}

type xmlBuilder struct {
	start xml.StartElement
	items []*xmlBuilder

	chardata interface{}
	cdata    bool // 表示 chardata 是一个 cdata 数据
}

func buildXML(p *ast.Request) ([]byte, error) {
	if p == nil || p.Type.V() == ast.TypeNone {
		return nil, nil
	}

	builder, err := parseXML(p.Param(), true, true)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	e := xml.NewEncoder(buf)
	e.Indent("", indent)

	if err = builder.encode(e); err != nil {
		return nil, err
	}

	if err := e.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseXML(p *ast.Param, chkArray, root bool) (*xmlBuilder, error) {
	builder := &xmlBuilder{
		start: xml.StartElement{
			Name: xml.Name{
				Local: p.Name.V(),
			},
			Attr: make([]xml.Attr, 0, len(p.Items)),
		},
		items: []*xmlBuilder{},
		cdata: p.XMLCData.V(),
	}
	if p.XMLNSPrefix != nil {
		builder.start.Name.Space = p.XMLNSPrefix.V()
	}

	if p.Array.V() && chkArray {
		if err := parseArray(p, builder); err != nil {
			return nil, err
		}
		if root {
			return builder.items[0], nil
		}
		return builder, nil
	}

	if p.Type.V() != ast.TypeObject {
		builder.chardata = genXMLValue(p)
		return builder, nil
	}

	for _, item := range p.Items {
		switch {
		case item.XMLAttr.V():
			attr := xml.Attr{
				Name: xml.Name{
					Local: item.Name.V(),
				},
				Value: fmt.Sprint(genXMLValue(item)),
			}
			if item.XMLNSPrefix != nil {
				attr.Name.Space = item.XMLNSPrefix.V()
			}
			builder.start.Attr = append(builder.start.Attr, attr)
		case item.XMLExtract.V():
			builder.chardata = genXMLValue(item)
			builder.cdata = item.XMLCData.V()
		case item.Array.V():
			if err := parseArray(item, builder); err != nil {
				return nil, err
			}
		default:
			b, err := parseXML(item, true, false)
			if err != nil {
				return nil, err
			}
			builder.items = append(builder.items, b)
		}
	} // end for

	return builder, nil
}

func parseArray(p *ast.Param, parent *xmlBuilder) error {
	b := parent
	if p.XMLWrapped.V() != "" {
		b = &xmlBuilder{items: []*xmlBuilder{}}
		if p.XMLNSPrefix != nil {
			b.start.Name.Space = p.XMLNSPrefix.V()
		}
		if p.XMLWrapped != nil {
			b.start.Name.Local = p.XMLWrapped.V()
		}
		parent.items = append(parent.items, b)
	}

	size := generateSliceSize()
	for i := 0; i < size; i++ {
		bb, err := parseXML(p, false, false)
		if err != nil {
			return err
		}
		b.items = append(b.items, bb)
	}

	return nil
}

func (builder *xmlBuilder) encode(e *xml.Encoder) error {
	if builder == nil {
		return nil
	}

	if builder.cdata && builder.chardata != nil {
		return e.EncodeElement(struct {
			string `xml:",cdata"`
		}{fmt.Sprint(builder.chardata)}, builder.start)
	} else if builder.chardata != nil {
		return e.EncodeElement(builder.chardata, builder.start)
	}

	if err := e.EncodeToken(builder.start); err != nil {
		return err
	}
	for _, item := range builder.items {
		if err := item.encode(e); err != nil {
			return err
		}
	}
	return e.EncodeToken(builder.start.End())
}

func genXMLValue(p *ast.Param) interface{} {
	switch p.Type.V() {
	case ast.TypeNone:
		return ""
	case ast.TypeBool:
		return generateBool()
	case ast.TypeNumber:
		return generateNumber(p)
	case ast.TypeString:
		return generateString(p)
	default: // ast.TypeObject:
		panic(fmt.Sprintf("无效的类型 %s", p.Type.V())) // 加载的时候已经作语法验证，此处还出错则直接 panic
	}
}
