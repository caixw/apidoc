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

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

type xmlValidator struct {
	param   *doc.Param
	decoder *xml.Decoder
	names   []string // 按顺序保存变量名称
}

func validXML(p *doc.Request, content []byte) error {
	if len(content) == 0 {
		if p == nil || p.Type == doc.None {
			return nil
		}
		return message.NewLocaleError("", "", 0, locale.ErrInvalidFormat)
	}

	validator := &xmlValidator{
		param:   p.ToParam(),
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
				if err := validator.validValue(attr.Name.Local, attr.Value); err != nil {
					return err
				}
			}
		case xml.EndElement:
			validator.popName()
		case xml.CharData:
			if err := validator.validValue("", string(v)); err != nil {
				return err
			}
		case xml.Comment, xml.ProcInst, xml.Directive:
		}
	}
}

// 如果 t == "" 表示不需要验证类型，比如 null 可以赋值给任何类型
func (validator *xmlValidator) validValue(k, v string) error {
	field := strings.Join(validator.names, "/")

	p := validator.find(k)
	if p == nil {
		if k != "" {
			field += "/" + k
		}
		return message.NewLocaleError("", field, 0, locale.ErrNotFound)
	}

	if p.XMLAttr {
		field += "/@" + k
	}

	switch p.Type {
	case doc.Number:
		if !is.Number(v) {
			return message.NewLocaleError("", field, 0, locale.ErrInvalidFormat)
		}
	case doc.Bool:
		if _, err := strconv.ParseBool(v); err != nil {
			return message.NewLocaleError("", field, 0, locale.ErrInvalidFormat)
		}
	case doc.String:
		return nil
	case doc.None:
		if v != "" {
			return message.NewLocaleError("", field, 0, locale.ErrInvalidValue)
		}
	case doc.Object:
	}

	if p.IsEnum() {
		for _, enum := range p.Enums {
			if enum.Value == v {
				return nil
			}
		}
		return message.NewLocaleError("", field, 0, locale.ErrInvalidValue)
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

// 如果 names 为空，返回 validator.param
//
// name 会被附加在 validator.names 之后，进行查询，如果为空，则只查询 validator.names 的值。
func (validator *xmlValidator) find(name string) *doc.Param {
	names := make([]string, len(validator.names))
	copy(names, validator.names)
	if len(name) > 0 {
		names = append(names, name)
	}

	if len(names) == 0 || validator.param == nil || names[0] != validator.param.Name {
		return nil
	}

	p := validator.param
LOOP:
	for i := 1; i < len(names); i++ {
		name := names[i]

		for _, pp := range p.Items {
			if pp.Name == name {
				p = pp
				continue LOOP
			}

			if pp.Array && pp.XMLWrapped == name {
				i--
				p = pp
				continue LOOP
			}
		}

		return nil
	}

	return p
}

type xmlBuilder struct {
	start    xml.StartElement
	charData string
	items    []*xmlBuilder
}

func buildXML(p *doc.Request) ([]byte, error) {
	if p == nil || p.Type == doc.None {
		return nil, nil
	}

	builder, err := parseXML(p.ToParam(), true, true)
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

func parseXML(p *doc.Param, chkArray, root bool) (*xmlBuilder, error) {
	if p == nil || p.Type == doc.None {
		return nil, nil
	}

	builder := &xmlBuilder{
		start: xml.StartElement{
			Name: xml.Name{
				Space: p.XMLNSPrefix,
				Local: p.Name,
			},
			Attr: make([]xml.Attr, 0, len(p.Items)),
		},
		items: []*xmlBuilder{},
	}

	if p.Array && chkArray {
		if err := parseArray(p, builder); err != nil {
			return nil, err
		}
		if root {
			return builder.items[0], nil
		}
		return builder, nil
	}

	if p.Type != doc.Object {
		switch p.Type {
		case doc.Bool:
			builder.charData = fmt.Sprint(generateBool())
		case doc.Number:
			builder.charData = fmt.Sprint(generateNumber(p))
		case doc.String:
			builder.charData = fmt.Sprint(generateString(p))
		}
		return builder, nil
	}

	for _, item := range p.Items {
		switch {
		case item.XMLAttr:
			v, err := getXMLValue(item)
			if err != nil {
				return nil, err
			}

			builder.start.Attr = append(builder.start.Attr, xml.Attr{
				Name: xml.Name{
					Space: item.XMLNSPrefix,
					Local: item.Name,
				},
				Value: fmt.Sprint(v),
			})
		case item.XMLExtract:
			builder.charData = fmt.Sprint(getXMLValue(item))
		case item.Array:
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

func parseArray(p *doc.Param, parent *xmlBuilder) error {
	b := parent
	if p.XMLWrapped != "" {
		b = &xmlBuilder{
			start: xml.StartElement{
				Name: xml.Name{
					Space: p.XMLNSPrefix,
					Local: p.XMLWrapped,
				},
			},
			items: []*xmlBuilder{},
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

	return e.EncodeToken(xml.EndElement{Name: builder.start.Name})
}

func getXMLValue(p *doc.Param) (interface{}, error) {
	switch p.Type {
	case doc.None:
		return "", nil
	case doc.Bool:
		return generateBool(), nil
	case doc.Number:
		return generateNumber(p), nil
	case doc.String:
		return generateString(p), nil
	default: // doc.Object:
		return nil, message.NewLocaleError("", "", 0, locale.ErrInvalidFormat)
	}
}
