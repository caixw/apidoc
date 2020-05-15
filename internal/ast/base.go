// SPDX-License-Identifier: MIT

package ast

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/issue9/version"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

type (
	// CData 表示 XML 的 CDATA 数据
	CData struct {
		token.BaseTag
		Value    token.String `apidoc:"-"`
		RootName struct{}     `apidoc:"string,meta,usage-string"`
	}

	// Content 表示一段普通的 XML 字符串
	Content struct {
		token.Base
		Value    string   `apidoc:"-"`
		RootName struct{} `apidoc:"string,meta,usage-string"`
	}

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
		token.BaseAttribute
		Value    token.String `apidoc:"-"`
		RootName struct{}     `apidoc:"string,meta,usage-string"`
	}

	// NumberAttribute 表示数值类型的属性
	NumberAttribute struct {
		token.BaseAttribute
		Value    Number   `apidoc:"-"`
		RootName struct{} `apidoc:"number,meta,usage-number"`
	}

	// BoolAttribute 表示布尔值类型的属性
	BoolAttribute struct {
		token.BaseAttribute
		Value    Bool     `apidoc:"-"`
		RootName struct{} `apidoc:"bool,meta,usage-bool"`
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
		token.BaseTag
		Content  Content  `apidoc:",content"`
		RootName struct{} `apidoc:"string,meta,usage-string"`
	}
)

// EncodeXML Encoder.EncodeXML
func (cdata *CData) EncodeXML() (string, error) {
	return cdata.Value.Value, nil
}

// EncodeXML Encoder.EncodeXML
func (s *Content) EncodeXML() (string, error) {
	return s.Value, nil
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *Attribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *Attribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// V 返回当前属性实际表示的值
func (a *Attribute) V() string {
	if a == nil {
		return ""
	}
	return a.Value.Value
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (num *NumberAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	v, err := strconv.Atoi(attr.Value.Value)
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.Value, err)
	}

	num.Value = Number{
		Range: attr.Value.Range,
		Value: v,
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (num *NumberAttribute) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(num.Value.Value), nil
}

// V 返回当前属性实际表示的值
func (num *NumberAttribute) V() int {
	if num == nil {
		return 0
	}
	return num.Value.Value
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (b *BoolAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	v, err := strconv.ParseBool(attr.Value.Value)
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.Value, err)
	}

	b.Value = Bool{
		Range: attr.Value.Range,
		Value: v,
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (b *BoolAttribute) EncodeXMLAttr() (string, error) {
	return strconv.FormatBool(b.Value.Value), nil
}

// V 返回当前属性实际表示的值
func (b *BoolAttribute) V() bool {
	if b == nil {
		return false
	}
	return b.Value.Value
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *MethodAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value
	a.Value.Value = strings.ToUpper(a.Value.Value)
	if !isValidMethod(a.Value.Value) {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.Value, locale.ErrInvalidValue)
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *MethodAttribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// V 返回当前属性实际表示的值
func (a *MethodAttribute) V() string {
	return (*Attribute)(a).V()
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *StatusAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	v := NumberAttribute{}
	if err := v.DecodeXMLAttr(p, attr); err != nil {
		return err
	}

	if !isValidStatus(v.Value.Value) {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.Value, locale.ErrInvalidValue)
	}

	*a = StatusAttribute(v)
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *StatusAttribute) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(a.Value.Value), nil
}

// V 返回当前属性实际表示的值
func (a *StatusAttribute) V() int {
	return (*NumberAttribute)(a).V()
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *TypeAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value
	if !isValidType(a.Value.Value) {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.Value, locale.ErrInvalidValue)
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *TypeAttribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// V 返回当前属性实际表示的值
func (a *TypeAttribute) V() string {
	return (*Attribute)(a).V()
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *VersionAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value
	if !isValidVersion(a.Value.Value) {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.Value, locale.ErrInvalidValue)
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *VersionAttribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// V 返回当前属性实际表示的值
func (a *VersionAttribute) V() string {
	return (*Attribute)(a).V()
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *APIDocVersionAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value

	ok, err := version.SemVerCompatible(Version, attr.Value.Value)
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.Value, err)
	}

	if !ok {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.Value, locale.ErrInvalidValue)
	}

	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *APIDocVersionAttribute) EncodeXMLAttr() (string, error) {
	return a.Value.Value, nil
}

// V 返回当前属性实际表示的值
func (a *APIDocVersionAttribute) V() string {
	return (*Attribute)(a).V()
}

// V 返回当前属性实际表示的值
func (s *Element) V() string {
	if s == nil {
		return ""
	}
	return s.Content.Value
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
