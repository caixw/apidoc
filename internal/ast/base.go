// SPDX-License-Identifier: MIT

package ast

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/issue9/version"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/token"
)

type (
	// CData 表示 XML 的 CDATA 数据
	CData = token.CData

	// String 表示一段普通的 XML 字符串
	String = token.String

	// Number 表示 XML 的数值类型
	Number struct {
		core.Range
		Value int
	}

	// Bool 表示 XML 的布尔值类型
	Bool struct {
		core.Range
		Value bool
	}

	// Attribute 表示 XML 属性
	Attribute struct {
		token.Base
		Value String
	}

	// NumberAttribute 表示数值类型的属性
	NumberAttribute struct {
		token.Base
		Value Number
	}

	// BoolAttribute 表示布尔值类型的属性
	BoolAttribute struct {
		token.Base
		Value Bool
	}

	// MethodAttribute 表示请求方法
	MethodAttribute Attribute

	// StatusAttribute 状态码的 XML 属性
	StatusAttribute NumberAttribute

	// TypeAttribute 表示方法类型属性
	TypeAttribute Attribute

	// VersionAttribute 表示版本号属性
	VersionAttribute Attribute

	// APIDocVersionAttribute 版本号属性，同时对版本号进行比较
	APIDocVersionAttribute Attribute

	// Element 定义不包含子元素和属性的基本的 XML 元素
	Element struct {
		token.Base
		Content String
	}
)

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *Attribute) DecodeXMLAttr(attr *token.Attribute) error {
	a.Value = attr.Value
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *Attribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (num *NumberAttribute) DecodeXMLAttr(attr *token.Attribute) error {
	v, err := strconv.Atoi(attr.Value.Value)
	if err != nil {
		return err
	}

	num.Value = Number{
		Range: attr.Value.Range,
		Value: v,
	}
	return err
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (num *NumberAttribute) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(num.Value.Value), nil
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (b *BoolAttribute) DecodeXMLAttr(attr *token.Attribute) error {
	v, err := strconv.ParseBool(attr.Value.Value)
	if err != nil {
		return err
	}

	b.Value = Bool{
		Range: attr.Value.Range,
		Value: v,
	}
	return err
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (b *BoolAttribute) EncodeXMLAttr() (string, error) {
	return strconv.FormatBool(b.Value.Value), nil
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *MethodAttribute) DecodeXMLAttr(attr *token.Attribute) error {
	a.Value = attr.Value
	if !isValidMethod(a.Value.Value) {
		return newError(attr.Value.Range, locale.ErrInvalidValue)
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *MethodAttribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *StatusAttribute) DecodeXMLAttr(attr *token.Attribute) error {
	v := NumberAttribute{}
	if err := v.DecodeXMLAttr(attr); err != nil {
		return err
	}

	if !isValidStatus(v.Value.Value) {
		return newError(attr.Value.Range, locale.ErrInvalidValue)
	}

	*a = StatusAttribute(v)
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *StatusAttribute) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(a.Value.Value), nil
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *TypeAttribute) DecodeXMLAttr(attr *token.Attribute) error {
	a.Value = attr.Value
	if !isValidType(a.Value.Value) {
		return newError(attr.Value.Range, locale.ErrInvalidValue)
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *TypeAttribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *VersionAttribute) DecodeXMLAttr(attr *token.Attribute) error {
	a.Value = attr.Value
	if !isValidVersion(a.Value.Value) {
		return newError(attr.Value.Range, locale.ErrInvalidValue)
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *VersionAttribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *APIDocVersionAttribute) DecodeXMLAttr(attr *token.Attribute) error {
	a.Value = attr.Value

	ok, err := version.SemVerCompatible(Version, attr.Value.Value)
	if err != nil {
		return withError(attr.Value.Range, err)
	}

	if !ok {
		return newError(attr.Value.Range, locale.ErrInvalidValue)
	}

	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *APIDocVersionAttribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// EncodeXML Encoder.EncodeXML
func (s *Element) EncodeXML() (string, error) {
	return s.Content.Value, nil
}

// DecodeXML Decoder.DecodeXML
func (s *Element) DecodeXML(p *token.Parser, start *token.StartElement) (*token.EndElement, error) {
	for {
		t, err := p.Token()
		if err == io.EOF {
			return nil, errors.New("test")
		} else if err != nil {
			return nil, err
		}

		switch elem := t.(type) {
		case *token.EndElement:
			if elem.Name.Value != start.Name.Value {
				return nil, p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
			}
			return elem, nil
		case *token.String:
			s.Content = *elem
		case *token.Instruction:
			return nil, p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
		}
	}
}

var validMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodHead,
	http.MethodOptions,
}

func isValidMethod(method string) bool {
	method = strings.ToUpper(method)
	for _, m := range validMethods {
		if m == method {
			return true
		}
	}

	return false
}

func isValidStatus(status int) bool {
	return (status >= http.StatusContinue) &&
		(status <= http.StatusNetworkAuthenticationRequired)
}

func isValidType(t string) bool {
	return t == TypeBool ||
		t == TypeObject ||
		t == TypeNumber ||
		t == TypeString ||
		t == TypeNone
}

func isValidVersion(v string) bool {
	return version.SemVerValid(v)
}
