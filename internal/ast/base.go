// SPDX-License-Identifier: MIT

package ast

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/issue9/version"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

const dateFormat = time.RFC3339

type (
	// CData 表示 XML 的 CDATA 数据
	CData struct {
		token.BaseTag
		Value    token.String `apidoc:"-"`
		RootName struct{}     `apidoc:"string,meta,usage-string"`
	}

	// Content 表示一段普通的 XML 元素内容
	Content struct {
		token.Base
		Value    string   `apidoc:"-"`
		RootName struct{} `apidoc:"string,meta,usage-string"`
	}

	// ExampleValue 示例代码的内容
	ExampleValue CData

	// Number 表示 XML 的数值类型
	Number struct {
		core.Range
		Int     int
		Float   float64
		IsFloat bool
	}

	// Bool 表示 XML 的布尔值类型
	Bool struct {
		core.Range
		Value bool
	}

	// Date 日期类型
	Date struct {
		core.Range
		Value time.Time
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

	// VersionAttribute 表示版本号属性
	VersionAttribute struct {
		token.BaseAttribute
		Value    token.String `apidoc:"-"`
		RootName struct{}     `apidoc:"version,meta,usage-version"`
	}

	// DateAttribute 日期属性
	DateAttribute struct {
		token.BaseAttribute
		Value    Date     `apidoc:"-"`
		RootName struct{} `apidoc:"date,meta,usage-date"`
	}

	// MethodAttribute 表示请求方法
	MethodAttribute Attribute

	// StatusAttribute 状态码的 XML 属性
	StatusAttribute NumberAttribute

	// TypeAttribute 表示方法类型属性
	TypeAttribute Attribute

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

// EncodeXML Encoder.EncodeXML
//
// 示例代码的内容，会在此处去掉其前导的空格
func (v *ExampleValue) EncodeXML() (string, error) {
	return trimLeftSpace(v.Value.Value), nil
}

func trimLeftSpace(v string) string {
	var min []byte // 找出的最小行首相同空格内容

	s := bufio.NewScanner(strings.NewReader(v))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Bytes()
		if len(bytes.TrimSpace(line)) == 0 { // 忽略空行
			continue
		}

		var index int
		for i, b := range line {
			if !unicode.IsSpace(rune(b)) {
				index = i
				break
			}
		}

		switch {
		case index == 0: // 当前行顶格
			return v
		case len(min) == 0: // 未初始化 min，且 index > 0
			min = make([]byte, index)
			copy(min, line[:index])
		default:
			min = getSamePrefix(min, line[:index])
		}
	}

	if len(min) == 0 {
		return v
	}

	buf := bufio.NewReader(strings.NewReader(v))
	ret := make([]byte, 0, buf.Size())
	for {
		line, err := buf.ReadBytes('\n')
		if bytes.HasPrefix(line, min) {
			line = line[len(min):]
		}
		ret = append(ret, line...)

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return string(ret)
}

func getSamePrefix(v1, v2 []byte) []byte {
	l1, l2 := len(v1), len(v2)
	l := l1
	if l1 > l2 {
		l = l2
	}

	for i := 0; i < l; i++ {
		if v1[i] != v2[i] {
			return v1[:i]
		}
	}
	return v1[:l]
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *Attribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
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
func (num *NumberAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	if v, err := strconv.Atoi(attr.Value.Value); err == nil {
		num.Value = Number{
			Range: attr.Value.Range,
			Int:   v,
		}
		return nil
	}

	v, err := strconv.ParseFloat(attr.Value.Value, 64)
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.String(), err)
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
func (b *BoolAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	v, err := strconv.ParseBool(attr.Value.Value)
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.String(), err)
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
func (a *MethodAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value
	a.Value.Value = strings.ToUpper(a.Value.Value)
	if !isValidMethod(a.Value.Value) {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.String(), locale.ErrInvalidValue)
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
func (a *StatusAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	v := NumberAttribute{}
	if err := v.DecodeXMLAttr(p, attr); err != nil {
		return err
	}

	if !isValidStatus(v.Value.Int) {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.String(), locale.ErrInvalidValue)
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
func (a *TypeAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value
	if !isValidType(a.Value.Value) {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.String(), locale.ErrInvalidValue)
	}
	return nil
}

// EncodeXMLAttr AttrEncoder.EncodeXMLAttr
func (a *TypeAttribute) EncodeXMLAttr() (string, error) {
	return a.V(), nil
}

// V 返回当前属性实际表示的值
func (a *TypeAttribute) V() string {
	return (*Attribute)(a).V()
}

// DecodeXMLAttr AttrDecoder.DecodeXMLAttr
func (a *VersionAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value
	if !isValidVersion(a.Value.Value) {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.String(), locale.ErrInvalidValue)
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
func (d *DateAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	t, err := time.Parse(dateFormat, attr.Value.Value)
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.String(), err)
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
func (a *APIDocVersionAttribute) DecodeXMLAttr(p *token.Parser, attr *token.Attribute) error {
	a.Value = attr.Value

	ok, err := version.SemVerCompatible(Version, attr.Value.Value)
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.String(), err)
	}

	if !ok {
		return p.NewError(attr.Value.Start, attr.Value.End, attr.Name.String(), locale.ErrInvalidValue)
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
