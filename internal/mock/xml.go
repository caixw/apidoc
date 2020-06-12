// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func validXML(ns []*ast.XMLNamespace, p *ast.Request, content []byte) error {
	if len(content) == 0 {
		if p == nil || p.Type.V() == ast.TypeNone {
			return nil
		}
		return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
	}

	d := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := d.Token()
		if errors.Is(err, io.EOF) && token == nil { // 正常结束
			return nil
		}
		if err != nil {
			return err
		}

		switch v := token.(type) {
		case xml.StartElement:
			if err := validXMLNamespaces(ns, v); err != nil {
				return err
			}
			return validXMLElement(ns, v, p.Param(), true, d, v.Name.Local)
		case xml.EndElement:
			return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
		}
	}
}

func validXMLNamespaces(ns []*ast.XMLNamespace, start xml.StartElement) error {
ATTR:
	for _, attr := range start.Attr {
		if attr.Name.Local == "xmlns" && attr.Name.Space == "" { // 默认
			for _, item := range ns {
				if item.URN.V() == attr.Value {
					continue ATTR
				}
			}
			return core.NewSyntaxError(core.Location{}, "xmlns", locale.ErrNotFound)
		} else if attr.Name.Space == "xmlns" && attr.Name.Local != "" {
			for _, item := range ns {
				if item.Prefix.V() == attr.Name.Local && item.URN.V() == attr.Value {
					continue ATTR
				}
			}
			return core.NewSyntaxError(core.Location{}, "xmlns:"+attr.Name.Local, locale.ErrNotFound)
		}
	}
	return nil
}

func validXMLName(name xml.Name, ns []*ast.XMLNamespace, p *ast.Param, chkArray bool) bool {
	if !p.Array.V() {
		if name.Local != p.Name.V() {
			return false
		} // else goto SPACE
	} else {
		n := parseXMLName(p, chkArray)
		if chkArray {
			if n == name.Local {
				goto SPACE
			}
			return n == ""
		}

		if name.Local != n {
			return false
		}
	}

SPACE:
	if name.Space == "" {
		return true
	}

	for _, item := range ns {
		if item.URN.V() == name.Space {
			return item.Prefix.V() == p.XMLNSPrefix.V()
		}
	}
	return false
}

func parseXMLName(p *ast.Param, parent bool) (name string) {
	if !p.Array.V() {
		return p.Name.V()
	}

	v := p.XMLWrapped.V()
	if v == "" {
		if parent {
			return ""
		}
		return p.Name.V()
	}

	index := strings.IndexByte(v, '>')
	switch {
	case index == 0:
		if parent {
			return ""
		}
		return v[1:]
	case index < 0:
		if parent {
			return v
		}
		return p.Name.V()
	default: // index > 0
		if parent {
			return v[:index]
		}
		return v[index+1:]
	}
}

func validXMLElement(
	ns []*ast.XMLNamespace,
	start xml.StartElement,
	p *ast.Param,
	chkArray bool, // 如果 p 为数组类型，则作为数组进行验证，否则只能当作普通元素进行验证
	d *xml.Decoder,
	field string,
) error {
	if err := validStartElement(ns, start, p, chkArray, field); err != nil {
		return err
	}

	for {
		token, err := d.Token()
		if errors.Is(err, io.EOF) && token == nil { // 正常结束
			return nil
		}
		if err != nil {
			return err
		}

		var chardata []byte
		switch v := token.(type) {
		case xml.StartElement:
			if chkArray && p.Array.V() && validXMLName(v.Name, ns, p, false) {
				return validXMLElement(ns, v, p, false, d, buildXMLField(field, p))
			}

			for _, pp := range p.Items {
				if validXMLName(v.Name, ns, pp, false) {
					if pp.XMLExtract.V() {
						return validXMLElement(ns, v, p, true, d, buildXMLField(field, pp))
					}
					return validXMLElement(ns, v, pp, true, d, buildXMLField(field, p))
				}
			}
		case xml.EndElement:
			if !validXMLName(v.Name, ns, p, chkArray) {
				return core.NewSyntaxError(core.Location{}, v.Name.Local, locale.ErrNotFoundEndTag)
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

func validStartElement(
	ns []*ast.XMLNamespace,
	start xml.StartElement,
	p *ast.Param,
	chkArray bool, // 如果 p 为数组类型，则作为数组进行验证，否则只能当作普通元素进行验证
	field string,
) error {
	for _, attr := range start.Attr {
		for _, pp := range p.Items {
			if !validXMLName(attr.Name, ns, pp, false) {
				continue
			}
			if err := validXMLParamValue(pp, buildXMLField(field, pp), attr.Value); err != nil {
				return err
			}
			break
		}
	}

	if p.Array.V() {
		if !validXMLName(start.Name, ns, p, chkArray) {
			return core.NewSyntaxError(core.Location{}, start.Name.Local, locale.ErrNotFound)
		}
	} else if !validXMLName(start.Name, ns, p, false) {
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

func buildXML(ns []*ast.XMLNamespace, p *ast.Request, indent string, g *GenOptions) ([]byte, error) {
	if p == nil || p.Type.V() == ast.TypeNone {
		return nil, nil
	}

	builder, err := parseXML(ns, p.Param(), true, true, g)
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

func parseXML(ns []*ast.XMLNamespace, p *ast.Param, chkArray, root bool, g *GenOptions) (*xmlBuilder, error) {
	builder := &xmlBuilder{
		start: xml.StartElement{
			Name: buildXMLName(p, chkArray),
			Attr: make([]xml.Attr, 0, len(p.Items)+len(ns)),
		},
		items: []*xmlBuilder{},
		cdata: p.XMLCData.V(),
	}

	if p.Array.V() && chkArray {
		if err := parseXMLArray(ns, p, builder, g); err != nil {
			return nil, err
		}
		if root {
			builder = builder.items[0]
		}
		goto RET
	}

	if p.Type.V() != ast.TypeObject {
		builder.chardata = genXMLValue(g, p)
		goto RET
	}

	for _, item := range p.Items {
		switch {
		case item.XMLAttr.V():
			attr := xml.Attr{
				Name:  buildXMLName(item, false),
				Value: fmt.Sprint(genXMLValue(g, item)),
			}
			builder.start.Attr = append(builder.start.Attr, attr)
		case item.XMLExtract.V():
			builder.chardata = genXMLValue(g, item)
			builder.cdata = item.XMLCData.V()
		case item.Array.V():
			if err := parseXMLArray(ns, item, builder, g); err != nil {
				return nil, err
			}
		default:
			b, err := parseXML(ns, item, true, false, g)
			if err != nil {
				return nil, err
			}
			builder.items = append(builder.items, b)
		}
	} // end for

RET:
	if root { // namespace
		for _, item := range ns {
			name := "xmlns"
			if item.Prefix.V() != "" {
				name += ":" + item.Prefix.V()
			}

			builder.start.Attr = append(builder.start.Attr, xml.Attr{
				Name:  xml.Name{Local: name},
				Value: item.URN.V(),
			})
		}
	}

	return builder, nil
}

func buildXMLName(p *ast.Param, chkArray bool) xml.Name {
	name := parseXMLName(p, chkArray)

	if p.XMLNSPrefix.V() != "" {
		return xml.Name{Local: p.XMLNSPrefix.V() + ":" + name}
	}

	return xml.Name{Local: name}
}

func parseXMLArray(ns []*ast.XMLNamespace, p *ast.Param, parent *xmlBuilder, g *GenOptions) error {
	b := parent
	if p.XMLWrapped.V() != "" && p.XMLWrapped.V()[0] != '>' {
		b = &xmlBuilder{items: []*xmlBuilder{}}
		b.start.Name = buildXMLName(p, true)
		parent.items = append(parent.items, b)
	}

	for i := 0; i < g.generateSliceSize(); i++ {
		bb, err := parseXML(ns, p, false, false, g)
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
		}{string: fmt.Sprint(builder.chardata)}, builder.start)
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

func genXMLValue(g *GenOptions, p *ast.Param) interface{} {
	switch p.Type.V() {
	case ast.TypeNone:
		return ""
	case ast.TypeBool:
		return g.generateBool()
	case ast.TypeNumber:
		return g.generateNumber(p)
	case ast.TypeString:
		return g.generateString(p)
	default: // ast.TypeObject:
		panic(fmt.Sprintf("无效的类型 %s", p.Type.V())) // 加载的时候已经作语法验证，此处还出错则直接 panic
	}
}
