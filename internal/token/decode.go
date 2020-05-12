// SPDX-License-Identifier: MIT

package token

import (
	"fmt"
	"io"
	"reflect"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Decoder 实现从 p 中解码内容到当前对象的值
type Decoder interface {
	// 从 p 中读取内容并实例化到当前对象中
	//
	// 必须要同时从 p 中读取相应的 EndElement 才能返回。
	// end 表示 EndElement.End 的值。
	//
	// NOTE: 如果是自闭合标签，则不会调用该接口。
	//
	// 接口应该只返回 *core.SyntaxError 作为错误对象。
	DecodeXML(p *Parser, start *StartElement) (end *EndElement, err error)
}

// AttrDecoder 实现从 attr 中解码内容到当前对象的值
type AttrDecoder interface {
	// 解析属性值
	//
	// 接口应该只返回 *core.SyntaxError 作为错误对象。
	DecodeXMLAttr(p *Parser, attr *Attribute) error
}

// Sanitizer 用于验证和修改对象中的数据
type Sanitizer interface {
	// 验证数据是否正确
	//
	// 可以通过 p.NewError 和 p.WithError 返回 *core.SyntaxError 类型的错误
	Sanitize(p *Parser) error
}

var (
	attrDecoderType = reflect.TypeOf((*AttrDecoder)(nil)).Elem()
	decoderType     = reflect.TypeOf((*Decoder)(nil)).Elem()
	sanitizerType   = reflect.TypeOf((*Sanitizer)(nil)).Elem()
)

// Decode 将 p 中的 XML 内容解码至 v 对象中
//
// Decode 中所有返回的错误对象，都可以转换成 *core.SyntaxError
func Decode(p *Parser, v interface{}) error {
	var hasRoot bool
	for {
		t, _, err := p.Token()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		switch elem := t.(type) {
		case *StartElement:
			if hasRoot { // 多个根元素
				return p.NewError(elem.Start, elem.End, "", locale.ErrInvalidXML)
			}
			hasRoot = true

			n := newNode("", reflect.ValueOf(v))
			if err := decodeElement(p, elem, n.value); err != nil {
				return err
			}
		case *EndElement:
			return p.NewError(elem.Start, elem.End, "", locale.ErrInvalidXML)
		}
	}
}

func (n *node) decode(p *Parser, start *StartElement) (*EndElement, error) {
	if err := n.decodeAttributes(p, start); err != nil {
		return nil, err
	}

	if start.Close {
		return nil, nil
	}

	end, err := n.decodeElements(p)
	if err != nil {
		return nil, err
	}

	// 判断 omitempty 属性
	for _, attr := range n.attrs {
		if attr.canNotEmpty() {
			return nil, p.NewError(start.Start, end.End, attr.name, locale.ErrRequired)
		}
	}
	for _, elem := range n.elems {
		if elem.canNotEmpty() {
			return nil, p.NewError(start.Start, end.End, elem.name, locale.ErrRequired)
		}
	}
	if n.cdata.canNotEmpty() {
		return nil, p.NewError(start.Start, end.End, "cdata", locale.ErrRequired)
	}
	if n.content.canNotEmpty() {
		return nil, p.NewError(start.Start, end.End, "content", locale.ErrRequired)
	}

	return end, nil
}

// 当前表示的值必须是一个非空值
func (v value) canNotEmpty() bool {
	return v.name != "" && // cdata 和 content 在未初始化时 name 字段为空值
		!v.omitempty &&
		(!v.IsValid() || !v.CanInterface() || is.Empty(v.Interface(), true))
}

// 将 start 的属性内容解码到 obj.attrs 之中
func (n *node) decodeAttributes(p *Parser, start *StartElement) error {
	if start == nil {
		return nil
	}

	for _, attr := range start.Attributes {
		item, found := n.attr(attr.Name.Value)
		if !found {
			continue
		}
		v := getRealValue(item.Value)
		v.Set(reflect.New(v.Type()).Elem())

		var impl bool
		if item.CanInterface() && item.Type().Implements(attrDecoderType) {
			if err := item.Interface().(AttrDecoder).DecodeXMLAttr(p, attr); err != nil {
				return err
			}
			impl = true
		} else if item.CanAddr() {
			pv := item.Addr()
			if pv.CanInterface() && pv.Type().Implements(attrDecoderType) {
				if err := pv.Interface().(AttrDecoder).DecodeXMLAttr(p, attr); err != nil {
					return err
				}
				impl = true
			}
		}

		if !impl {
			panic(fmt.Sprintf("当前属性 %s 未实现 AttrDecoder 接口", attr.Name.Value))
		}

		if err := setValue(item.Value, item.usage, p, attr.Start, attr.End, attr.Name, String{}); err != nil {
			return err
		}
	}

	return nil
}

func (n *node) decodeElements(p *Parser) (*EndElement, error) {
	for {
		t, _, err := p.Token()
		if err == io.EOF {
			// 应该只有 EndElement 才能返回，否则就不完整的 XML
			return nil, p.NewError(p.Position().Position, p.Position().Position, "", locale.ErrInvalidXML)
		} else if err != nil {
			return nil, err
		}

		switch elem := t.(type) {
		case *EndElement: // 找到当前对象的结束标签
			if elem.Name.Value == n.value.name {
				return elem, nil
			}
			return nil, p.NewError(elem.Start, elem.End, n.value.name, locale.ErrNotFoundEndTag)
		case *CData:
			if n.cdata.IsValid() {
				getRealValue(n.cdata.Value).Set(getRealValue(reflect.ValueOf(elem)))
			}
		case *String:
			if n.content.IsValid() {
				getRealValue(n.content.Value).Set(getRealValue(reflect.ValueOf(elem)))
			}
		case *StartElement:
			item, found := n.elem(elem.Name.Value)
			if !found {
				panic(fmt.Sprintf("不存在的子元素 %s", elem.Name.Value))
			}
			if err = decodeElement(p, elem, item); err != nil {
				return nil, err
			}
		}
	}
}

func decodeElement(p *Parser, start *StartElement, v value) error {
	v.Value = getRealValue(v.Value)
	k := v.Kind()
	switch {
	case k == reflect.Ptr, k == reflect.Func, k == reflect.Chan, k == reflect.Array, isPrimitive(v.Value):
		panic(fmt.Sprintf("%s 是无效的类型", v.Value.Type()))
	case k == reflect.Slice:
		return decodeSlice(p, start, v)
	}

	end, impl, err := callDecodeXML(v.Value, p, start)
	if !impl {
		end, err = newNode(start.Name.Value, v.Value).decode(p, start)
	}
	if err != nil {
		return err
	}
	return setElementValue(v.Value, v.usage, p, start, end)
}

func decodeSlice(p *Parser, start *StartElement, slice value) (err error) {
	// 不相配，表示当前元素找不到与之相配的元素，需要忽略这个元素，
	// 所以要过滤与 start 想匹配的结束符号才算结束。
	if !start.Close && (start.Name.Value != slice.name) {
		return findEndElement(p, start)
	}

	elem := reflect.New(slice.Type().Elem()).Elem()
	if elem.Kind() == reflect.Ptr && elem.IsNil() {
		elem.Set(reflect.New(elem.Type().Elem()))
	}

	end, impl, err := callDecodeXML(elem, p, start)
	if !impl {
		if isPrimitive(elem) {
			panic(fmt.Sprintf("%s:%s 必须实现 Decoder 接口", slice.name, elem.Type()))
		}
		end, err = newNode(start.Name.Value, elem).decode(p, start)
	}
	if err != nil {
		return err
	}

	if err = setElementValue(elem, slice.usage, p, start, end); err != nil {
		return err
	}
	slice.Value.Set(reflect.Append(slice.Value, elem))
	return nil
}

// 调用 v 的 DecodeXML 接口方法
func callDecodeXML(v reflect.Value, p *Parser, start *StartElement) (end *EndElement, impl bool, err error) {
	if v.CanInterface() && v.Type().Implements(decoderType) {
		if !start.Close {
			end, err = v.Interface().(Decoder).DecodeXML(p, start)
		}
		return end, true, err
	} else if v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(decoderType) {
			if !start.Close {
				end, err = pv.Interface().(Decoder).DecodeXML(p, start)
			}
			return end, true, err
		}
	}
	return nil, false, nil
}

// 找到与 start 相对应的结束符号位置
func findEndElement(p *Parser, start *StartElement) error {
	level := 0
	for {
		t, _, err := p.Token()
		if err == io.EOF {
			return p.NewError(start.Start, start.End, start.Name.Value, locale.ErrNotFoundEndTag) // 找不到相配的结束符号
		} else if err != nil {
			return err
		}

		switch elem := t.(type) {
		case *StartElement:
			if elem.Name.Value == start.Name.Value {
				level++
			}
		case *EndElement:
			if level == 0 && (elem.Name.Value == start.Name.Value) {
				return nil
			}
			level--
		}
	}
}

func setElementValue(v reflect.Value, usage string, p *Parser, start *StartElement, end *EndElement) error {
	if end == nil {
		return setValue(v, usage, p, start.Start, start.End, start.Name, String{})
	}
	return setValue(v, usage, p, start.Start, end.End, start.Name, end.Name)
}

func setValue(v reflect.Value, usage string, p *Parser, start, end core.Position, xmlName, xmlNameEnd String) error {
	if v.CanInterface() && v.Type().Implements(sanitizerType) {
		if err := v.Interface().(Sanitizer).Sanitize(p); err != nil {
			return err
		}
	} else if v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(sanitizerType) {
			if err := pv.Interface().(Sanitizer).Sanitize(p); err != nil {
				return err
			}
		}
	}

	v = getRealValue(v)
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("无效的 kind 类型: %s:%s", v.Type(), v.Kind()))
	}

	v.FieldByName(rangeName).Set(reflect.ValueOf(core.Range{Start: start, End: end}))

	if usage != "" { // CDATA 和 content 节点类型的 usage 内容为空
		v.FieldByName(usageKeyName).Set(reflect.ValueOf(usage))
	}

	if xmlName.Value != "" { // CDATA 和 content 节点的 StartTag 肯定为空
		v.FieldByName(elementTagName).Set(reflect.ValueOf(xmlName))
	}

	if xmlNameEnd.Value != "" {
		v.FieldByName(elementTagEndName).Set(reflect.ValueOf(xmlNameEnd))
	}

	return nil
}
