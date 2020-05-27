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
			return validElement(v, p.Param(), true, d)
		case xml.EndElement:
			return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
		}
	}
}

func validElement(start xml.StartElement, p *ast.Param, allowArray bool, decoder *xml.Decoder) error {
	for _, attr := range start.Attr {
		for _, pp := range p.Items {
			if attr.Name.Local != pp.Name.V() {
				continue
			}
			if err := validXMLParamValue(pp, pp.Name.V(), attr.Value); err != nil {
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
		return validElement(start, p, false, decoder)
	}

	if (p.Name.V() != start.Name.Local) && (!allowArray && p.XMLWrapped.V() != start.Name.Local) {
		return core.NewSyntaxError(core.Location{}, start.Name.Local, locale.ErrNotFound)
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF && token == nil { // 正常结束
			return nil
		}
		if err != nil {
			return err
		}

		var chardata []byte
		switch v := token.(type) {
		case xml.StartElement:
			chardata = nil // 有子元素，说明 chardata 无用

			if allowArray && p.Array.V() && v.Name.Local == p.Name.V() {
				return validElement(v, p, false, decoder)
			}

			for _, pp := range p.Items {
				if pp.Name.V() == v.Name.Local {
					if pp.XMLExtract.V() {
						return validElement(v, p, true, decoder)
					}
					return validElement(v, pp, true, decoder)
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
		case xml.Comment, xml.ProcInst, xml.Directive:
		}
	}
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
	start    xml.StartElement
	charData string
	items    []*xmlBuilder
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
		switch p.Type.V() {
		case ast.TypeBool:
			builder.charData = fmt.Sprint(generateBool())
		case ast.TypeNumber:
			builder.charData = fmt.Sprint(generateNumber(p))
		case ast.TypeString:
			builder.charData = fmt.Sprint(generateString(p))
		}
		return builder, nil
	}

	for _, item := range p.Items {
		switch {
		case item.XMLAttr.V():
			v, err := getXMLValue(item)
			if err != nil {
				return nil, err
			}

			attr := xml.Attr{
				Name: xml.Name{
					Local: item.Name.V(),
				},
				Value: fmt.Sprint(v),
			}
			if item.XMLNSPrefix != nil {
				attr.Name.Space = item.XMLNSPrefix.V()
			}
			builder.start.Attr = append(builder.start.Attr, attr)
		case item.XMLExtract.V():
			v, err := getXMLValue(item)
			if err != nil {
				return nil, err
			}

			builder.charData = fmt.Sprint(v)
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

	if builder.charData != "" {
		return e.EncodeElement(builder.charData, builder.start)
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

func getXMLValue(p *ast.Param) (interface{}, error) {
	switch p.Type.V() {
	case ast.TypeNone:
		return "", nil
	case ast.TypeBool:
		return generateBool(), nil
	case ast.TypeNumber:
		return generateNumber(p), nil
	case ast.TypeString:
		return generateString(p), nil
	default: // ast.TypeObject:
		return nil, core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
	}
}
