// SPDX-License-Identifier: MIT

package ast

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/issue9/version"

	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

const dateFormat = time.RFC3339

type (
	// Attribute 表示 XML 属性
	Attribute struct {
		xmlenc.BaseAttribute
		Value    xmlenc.String `apidoc:"-"`
		RootName struct{}      `apidoc:"string,meta,usage-string"`
	}

	// NumberAttribute 表示数值类型的属性
	NumberAttribute struct {
		xmlenc.BaseAttribute
		Value    Number   `apidoc:"-"`
		RootName struct{} `apidoc:"number,meta,usage-number"`
	}

	// BoolAttribute 表示布尔值类型的属性
	BoolAttribute struct {
		xmlenc.BaseAttribute
		Value    Bool     `apidoc:"-"`
		RootName struct{} `apidoc:"bool,meta,usage-bool"`
	}

	// VersionAttribute 表示版本号属性
	VersionAttribute struct {
		xmlenc.BaseAttribute
		Value    xmlenc.String `apidoc:"-"`
		RootName struct{}      `apidoc:"version,meta,usage-version"`
	}

	// DateAttribute 日期属性
	DateAttribute struct {
		xmlenc.BaseAttribute
		Value    Date     `apidoc:"-"`
		RootName struct{} `apidoc:"date,meta,usage-date"`
	}

	// MethodAttribute 表示请求方法
	MethodAttribute Attribute

	// StatusAttribute 状态码的 XML 属性
	StatusAttribute NumberAttribute

	// TypeAttribute 表示方法类型属性
	TypeAttribute struct {
		xmlenc.BaseAttribute
		Value    xmlenc.String `apidoc:"-"`
		RootName struct{}      `apidoc:"type,meta,usage-type"`
	}

	// APIDocVersionAttribute 版本号属性，同时对版本号进行比较
	APIDocVersionAttribute Attribute
)

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *Attribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	a.Value = attr.Value
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *Attribute) EncodeXMLAttr() (string, error) {
	return a.V(), nil
}

// V 返回当前属性实际表示的值
func (a *Attribute) V() string {
	if a == nil {
		return ""
	}
	return a.Value.Value
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (num *NumberAttribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	if v, err := strconv.Atoi(attr.Value.Value); err == nil {
		num.Value = Number{
			Range: attr.Value.Range,
			Int:   v,
		}
		return nil
	}

	v, err := strconv.ParseFloat(attr.Value.Value, 64)
	if err != nil {
		return attr.Value.WithError(err).WithField(attr.Name.String())
	}
	num.Value = Number{
		Range:   attr.Value.Range,
		Float:   v,
		IsFloat: true,
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (num *NumberAttribute) EncodeXMLAttr() (string, error) {
	if num.IsFloat() {
		return strconv.FormatFloat(num.FloatValue(), 'f', -1, 64), nil
	}
	return strconv.Itoa(num.IntValue()), nil
}

// IntValue 返回当前属性实际表示的值
func (num *NumberAttribute) IntValue() int {
	if num == nil {
		return 0
	}
	return num.Value.Int
}

// FloatValue 返回当前属性实际表示的值
func (num *NumberAttribute) FloatValue() float64 {
	if num == nil {
		return 0
	}
	return num.Value.Float
}

// IsFloat 当前的数值类型是否为浮点型
func (num *NumberAttribute) IsFloat() bool {
	return num.Value.IsFloat
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (b *BoolAttribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	v, err := strconv.ParseBool(attr.Value.Value)
	if err != nil {
		return attr.Value.WithError(err).WithField(attr.Name.String())
	}

	b.Value = Bool{
		Range: attr.Value.Range,
		Value: v,
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (b *BoolAttribute) EncodeXMLAttr() (string, error) {
	return strconv.FormatBool(b.V()), nil
}

// V 返回当前属性实际表示的值
func (b *BoolAttribute) V() bool {
	if b == nil {
		return false
	}
	return b.Value.Value
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *MethodAttribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	a.Value = attr.Value
	a.Value.Value = strings.ToUpper(a.V())
	if !isValidMethod(a.V()) {
		return attr.Value.NewError(locale.ErrInvalidValue).WithField(attr.Name.String())
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *MethodAttribute) EncodeXMLAttr() (string, error) {
	return a.V(), nil
}

// V 返回当前属性实际表示的值
func (a *MethodAttribute) V() string {
	return (*Attribute)(a).V()
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *StatusAttribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	v := NumberAttribute{}
	if err := v.DecodeXMLAttr(p, attr); err != nil {
		return err
	}

	if !isValidStatus(v.Value.Int) {
		return attr.Value.NewError(locale.ErrInvalidValue).WithField(attr.Name.String())
	}

	*a = StatusAttribute(v)
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *StatusAttribute) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(a.V()), nil
}

// V 返回当前属性实际表示的值
func (a *StatusAttribute) V() int {
	return (*NumberAttribute)(a).IntValue()
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *TypeAttribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	a.Value = attr.Value
	if !isValidType(a.V()) {
		return attr.Value.NewError(locale.ErrInvalidValue).WithField(attr.Name.String())
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *TypeAttribute) EncodeXMLAttr() (string, error) {
	return a.V(), nil
}

// V 返回当前属性实际表示的值
func (a *TypeAttribute) V() string {
	if a == nil {
		return ""
	}
	return a.Value.Value
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *VersionAttribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	a.Value = attr.Value
	if !isValidVersion(a.V()) {
		return attr.Value.NewError(locale.ErrInvalidValue).WithField(attr.Name.String())
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *VersionAttribute) EncodeXMLAttr() (string, error) {
	return a.V(), nil
}

// V 返回当前属性实际表示的值
func (a *VersionAttribute) V() string {
	return a.Value.Value
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (d *DateAttribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	t, err := time.Parse(dateFormat, attr.Value.Value)
	if err != nil {
		return attr.Value.WithError(err).WithField(attr.Name.String())
	}

	d.Value = Date{
		Range: attr.Value.Range,
		Value: t,
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (d *DateAttribute) EncodeXMLAttr() (string, error) {
	return d.V().Format(dateFormat), nil
}

// V 返回当前属性实际表示的值
func (d *DateAttribute) V() time.Time {
	return d.Value.Value
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *APIDocVersionAttribute) DecodeXMLAttr(p *xmlenc.Parser, attr *xmlenc.Attribute) error {
	a.Value = attr.Value

	ok, err := version.SemVerCompatible(Version, attr.Value.Value)
	if err != nil {
		return attr.Value.WithError(err).WithField(attr.Name.String())
	}

	if !ok {
		return attr.Value.NewError(locale.ErrInvalidValue).WithField(attr.Name.String())
	}

	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *APIDocVersionAttribute) EncodeXMLAttr() (string, error) {
	return a.V(), nil
}

// V 返回当前属性实际表示的值
func (a *APIDocVersionAttribute) V() string {
	return (*Attribute)(a).V()
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
		t == TypeInt ||
		t == TypeFloat ||
		t == TypeString ||
		t == TypeURL ||
		t == TypeEmail ||
		t == TypeImage ||
		t == TypeDate ||
		t == TypeTime ||
		t == TypeDateTime ||
		t == TypeNone
}

func isValidVersion(v string) bool {
	return version.SemVerValid(v)
}
