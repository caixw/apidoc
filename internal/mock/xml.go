// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/ast"
	"github.com/caixw/apidoc/v6/internal/locale"
)

type xmlValidator struct {
	param   *ast.Param
	decoder *xml.Decoder
	names   []string // 按顺序保存变量名称
}

func validXML(p *ast.Request, content []byte) error {
	if len(content) == 0 {
		if p == nil || p.Type == nil || p.Type.Value.Value == ast.TypeNone {
			return nil
		}
		return core.NewLocaleError(core.Location{}, "", locale.ErrInvalidFormat)
	}

	validator := &xmlValidator{
		param:   p.Param(),
		decoder: xml.NewDecoder(bytes.NewReader(content)),
		names:   []string{},
	}

	return validator.valid()
}

func (validator *xmlValidator) valid() error {
	for {
		token, err := validator.decoder.Token()
		if err == io.EOF && token == nil { // 正常结束
			return nil
		}

		if err != nil {
			return err
		}

		switch v := token.(type) {
		case xml.StartElement:
			validator.pushName(v.Name.Local)
			for _, attr := range v.Attr {
				validator.pushName(attr.Name.Local)
				if err := validator.validValue(attr.Value); err != nil {
					return err
				}
				validator.popName()
			}
		case xml.EndElement:
			validator.popName()
		case xml.CharData:
			if len(v) > 0 && v[0] == '\n' && (len(v[1:])%len(indent) == 0) {
				continue
			}

			if err := validator.validValue(string(v)); err != nil {
				return err
			}
		case xml.Comment, xml.ProcInst, xml.Directive:
		}
	}
}

// 如果 t == "" 表示不需要验证类型，比如 null 可以赋值给任何类型
func (validator *xmlValidator) validValue(v string) error {
	field := strings.Join(validator.names, "/")

	p := validator.find()
	if p == nil {
		return core.NewLocaleError(core.Location{}, field, locale.ErrNotFound)
	}

	return validXMLParamValue(p, field, v)
}

// 验证 p 描述的类型与 v 是否匹配，如果不匹配返回错误信息。
// field 表示 p 在整个对象中的位置信息。
func validXMLParamValue(p *ast.Param, field, v string) error {
	switch {
	case p.Type == nil || p.Type.Value.Value == ast.TypeNone:
		if v != "" {
			return core.NewLocaleError(core.Location{}, field, locale.ErrInvalidValue)
		}
	case p.Type.Value.Value == ast.TypeNumber:
		if !is.Number(v) {
			return core.NewLocaleError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case p.Type.Value.Value == ast.TypeBool:
		if _, err := strconv.ParseBool(v); err != nil {
			return core.NewLocaleError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case p.Type.Value.Value == ast.TypeString:
		return nil
	default: // case doc.Object:
		return core.NewLocaleError(core.Location{}, field, locale.ErrInvalidFormat)
	}

	if isEnum(p) {
		for _, enum := range p.Enums {
			if enum.Value.Value.Value == v {
				return nil
			}
		}
		return core.NewLocaleError(core.Location{}, field, locale.ErrInvalidValue)
	}

	return nil
}

func (validator *xmlValidator) pushName(name string) *xmlValidator {
	validator.names = append(validator.names, name)
	return validator
}

func (validator *xmlValidator) popName() *xmlValidator {
	if len(validator.names) > 0 {
		validator.names = validator.names[:len(validator.names)-1]
	}
	return validator
}

// 如果 names 为空，返回 nil
func (validator *xmlValidator) find() *ast.Param {
	p := validator.param

	if len(validator.names) == 0 || p == nil {
		return nil
	}

	var start int
	if (p.Array != nil && p.Array.Value.Value) &&
		(p.XMLWrapped != nil && p.XMLWrapped.Value.Value == validator.names[0]) {
		if len(validator.names) > 1 && validator.names[1] == p.Name.Value.Value {
			start = 2
		} else {
			return nil
		}
	} else if p.Name != nil && (p.Name.Value.Value == validator.names[0]) {
		start = 1
	} else {
		return nil
	}

	names := validator.names[start:]

LOOP:
	for i := 0; i < len(names); i++ {
		name := names[i]

		for _, pp := range p.Items {
			if (pp.Array != nil && pp.Array.Value.Value) &&
				(pp.XMLWrapped != nil && pp.XMLWrapped.Value.Value == name) {
				i++
				if i < len(names) && pp.Name.Value.Value == names[i] {
					p = pp
					continue LOOP
				}
				return nil
			}

			if pp.Name.Value.Value == name {
				p = pp
				continue LOOP
			}
		}

		return nil
	}

	// 如果根据 names 查找出来的实例带 XMLExtract，则肯定有问题。
	if p.XMLExtract != nil && p.XMLExtract.Value.Value {
		return nil
	}

	// 从子项中查找带 XMLExtract 的项
	if p.Type.Value.Value == ast.TypeObject {
		for _, pp := range p.Items {
			if pp.XMLExtract != nil && pp.XMLExtract.Value.Value {
				p = pp
				break
			}
		}
	}

	return p
}

type xmlBuilder struct {
	start    xml.StartElement
	charData string
	items    []*xmlBuilder
}

func buildXML(p *ast.Request) ([]byte, error) {
	if p == nil || p.Type == nil || p.Type.Value.Value == ast.TypeNone {
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
				Local: p.Name.Value.Value,
			},
			Attr: make([]xml.Attr, 0, len(p.Items)),
		},
		items: []*xmlBuilder{},
	}
	if p.XMLNSPrefix != nil {
		builder.start.Name.Space = p.XMLNSPrefix.Value.Value
	}

	if (p.Array != nil && p.Array.Value.Value) && chkArray {
		if err := parseArray(p, builder); err != nil {
			return nil, err
		}
		if root {
			return builder.items[0], nil
		}
		return builder, nil
	}

	if p.Type != nil && p.Type.Value.Value != ast.TypeObject {
		switch p.Type.Value.Value {
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
		case item.XMLAttr != nil && item.XMLAttr.Value.Value:
			v, err := getXMLValue(item)
			if err != nil {
				return nil, err
			}

			attr := xml.Attr{
				Name: xml.Name{
					Local: item.Name.Value.Value,
				},
				Value: fmt.Sprint(v),
			}
			if item.XMLNSPrefix != nil {
				attr.Name.Space = item.XMLNSPrefix.Value.Value
			}
			builder.start.Attr = append(builder.start.Attr, attr)
		case item.XMLExtract != nil && item.XMLExtract.Value.Value:
			v, err := getXMLValue(item)
			if err != nil {
				return nil, err
			}

			builder.charData = fmt.Sprint(v)
		case item.Array != nil && item.Array.Value.Value:
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
	if p.XMLWrapped != nil && p.XMLWrapped.Value.Value != "" {
		b = &xmlBuilder{items: []*xmlBuilder{}}
		if p.XMLNSPrefix != nil {
			b.start.Name.Space = p.XMLNSPrefix.Value.Value
		}
		if p.XMLWrapped != nil {
			b.start.Name.Local = p.XMLWrapped.Value.Value
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
	switch {
	case p.Type == nil || p.Type.Value.Value == ast.TypeNone:
		return "", nil
	case p.Type.Value.Value == ast.TypeBool:
		return generateBool(), nil
	case p.Type.Value.Value == ast.TypeNumber:
		return generateNumber(p), nil
	case p.Type.Value.Value == ast.TypeString:
		return generateString(p), nil
	default: // doc.Object:
		return nil, core.NewLocaleError(core.Location{}, "", locale.ErrInvalidFormat)
	}
}
