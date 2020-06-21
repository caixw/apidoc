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

type xmlValidator struct {
	namespaces []*ast.XMLNamespace
	decoder    *xml.Decoder
}

func validXML(ns []*ast.XMLNamespace, p *ast.Request, content []byte) error {
	if len(content) == 0 {
		if p == nil || p.Type.V() == ast.TypeNone {
			return nil
		}
		return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
	}

	validator := &xmlValidator{
		namespaces: ns,
		decoder:    xml.NewDecoder(bytes.NewReader(content)),
	}
	for {
		token, err := validator.decoder.Token()
		if errors.Is(err, io.EOF) && token == nil { // 正常结束
			return nil
		}
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			if err := validator.validXMLNamespaces(elem); err != nil {
				return err
			}
			return validator.validXMLElement(elem, p.Param(), true, elem.Name.Local)
		case xml.EndElement:
			return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
		}
	}
}

func (v *xmlValidator) validXMLNamespaces(start xml.StartElement) error {
ATTR:
	for _, attr := range start.Attr {
		if attr.Name.Local == "xmlns" && attr.Name.Space == "" { // 默认
			for _, item := range v.namespaces {
				if item.URN.V() == attr.Value {
					continue ATTR
				}
			}
			return core.NewSyntaxError(core.Location{}, "xmlns", locale.ErrNotFound)
		} else if attr.Name.Space == "xmlns" && attr.Name.Local != "" {
			for _, item := range v.namespaces {
				if item.Prefix.V() == attr.Name.Local && item.URN.V() == attr.Value {
					continue ATTR
				}
			}
			return core.NewSyntaxError(core.Location{}, "xmlns:"+attr.Name.Local, locale.ErrNotFound)
		}
	}
	return nil
}

func (v *xmlValidator) validXMLName(name xml.Name, p *ast.Param, chkArray bool) bool {
	if !p.Array.V() {
		if name.Local != p.Name.V() {
			return false
		} // else goto SPACE
	} else {
		n := parseXMLWrappedName(p, chkArray)
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

	for _, item := range v.namespaces {
		if item.URN.V() == name.Space {
			return item.Prefix.V() == p.XMLNSPrefix.V()
		}
	}
	return false
}

// parent 如果是数组，则是否拿 wrapped 中指示的父元素名称。
func parseXMLWrappedName(p *ast.Param, parent bool) (name string) {
	if !p.Array.V() {
		return p.Name.V()
	}

	choose := func(v1, v2 string) string {
		if parent {
			return v1
		}
		return v2
	}

	v := p.XMLWrapped.V()
	if v == "" {
		return choose("", p.Name.V())
	}

	index := strings.IndexByte(v, '>')
	switch {
	case index == 0:
		return choose("", v[1:])
	case index < 0:
		return choose(v, p.Name.V())
	default: // index > 0
		return choose(v[:index], v[index+1:])
	}
}

func (v *xmlValidator) validXMLElement(start xml.StartElement, p *ast.Param, chkArray bool, field string) error {
	if err := v.validStartElement(start, p, chkArray, field); err != nil {
		return err
	}

	var chardata []byte
	var started bool
LOOP:
	for {
		token, err := v.decoder.Token()
		if errors.Is(err, io.EOF) && token == nil { // 正常结束
			return nil
		}
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:
			if chkArray && p.Array.V() && v.validXMLName(elem.Name, p, false) {
				if err := v.validXMLElement(elem, p, false, buildXMLField(field, p)); err != nil {
					return err
				}
				chardata = nil
				started = true
				continue LOOP
			}

			for _, pp := range p.Items {
				if v.validXMLName(elem.Name, pp, true) {
					if pp.XMLExtract.V() {
						pp = p
					}
					if err = v.validXMLElement(elem, pp, true, buildXMLField(field, pp)); err != nil {
						return err
					}
					chardata = nil
					started = true
					continue LOOP
				}
			}
		case xml.EndElement:
			if !v.validXMLName(elem.Name, p, chkArray) {
				return core.NewSyntaxError(core.Location{}, elem.Name.Local, locale.ErrNotFoundEndTag)
			}

			if chardata != nil && !started {
				return validXMLValue(p, p.Name.V(), string(chardata))
			}
			return nil
		case xml.CharData:
			if !started {
				chardata = elem.Copy()
			}
		}
	}
}

func (v *xmlValidator) validStartElement(start xml.StartElement, p *ast.Param, chkArray bool, field string) error {
	for _, attr := range start.Attr {
		for _, pp := range p.Items {
			if !v.validXMLName(attr.Name, pp, false) {
				continue
			}
			if err := validXMLValue(pp, buildXMLField(field, pp), attr.Value); err != nil {
				return err
			}
			break
		}
	}

	if !v.validXMLName(start.Name, p, chkArray) {
		return core.NewSyntaxError(core.Location{}, start.Name.Local, locale.ErrNotFound)
	}
	return nil
}

func buildXMLField(field string, p *ast.Param) string {
	if p.XMLAttr.V() {
		return field + "@" + p.Name.V()
	}
	return field + ">" + p.Name.V()
}

// 验证 p 描述的类型与 v 是否匹配，如果不匹配返回错误信息。
// field 表示 p 在整个对象中的位置信息。
func validXMLValue(p *ast.Param, field, v string) error {
	switch p.Type.V() {
	case ast.TypeNone:
		if v != "" {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidValue)
		}
	case ast.TypeNumber:
		if !is.Number(v) {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeInt:
		if _, err := strconv.ParseInt(v, 10, 64); err != nil {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeFloat:
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeBool:
		if _, err := strconv.ParseBool(v); err != nil {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeEmail:
		if !is.Email(v) {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeURL:
		if !is.URL(v) {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeString, ast.TypeObject, ast.TypeImage:
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
	name := parseXMLWrappedName(p, chkArray)

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
	switch primitive, _ := ast.ParseType(p.Type.V()); primitive {
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
