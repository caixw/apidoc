// SPDX-License-Identifier: MIT

package token

import (
	"fmt"
	"reflect"

	"github.com/issue9/conv"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
)

// Decoder 实现从 p 中解码内容到当前对象
type Decoder interface {
	// 从 p 中读取内容并实例化到当前对象中
	//
	// 必须要同时从 p 中读取相应的 EndElement 才能返回
	DecodeXML(p *Parser, start *StartElement) error
}

// AttrDecoder 实现从 attr 中解码内容到当前对象
type AttrDecoder interface {
	DecodeXMLAttr(attr *Attribute) error
}

var (
	attrDecoderType = reflect.TypeOf((*AttrDecoder)(nil)).Elem()
	decoderType     = reflect.TypeOf((*Decoder)(nil)).Elem()
)

// Decode 将 b 中的 XML 内容解码至 v 对象中
func Decode(b core.Block, v interface{}) error {
	p, err := NewParser(b)
	if err != nil {
		return err
	}

	var hasRoot bool
	for {
		t, err := p.Token()
		if err != nil {
			return err
		}
		if t == nil {
			return nil
		}

		switch elem := t.(type) {
		case *StartElement:
			if hasRoot { // 多个根元素
				return p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
			}

			o := newNode(elem.Name.Value, reflect.ValueOf(v))
			if err := o.decode(p, elem); err != nil {
				return err
			}
		case *EndElement:
			return p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
		}
	}
}

func (n *node) decode(p *Parser, start *StartElement) error {
	if err := n.decodeAttributes(start); err != nil {
		return err
	}

	if start.Close {
		return nil
	}

	return n.decodeElements(p)
}

// 将 start 的属性内容解码到 obj.attrs 之中
func (n *node) decodeAttributes(start *StartElement) error {
	if start == nil || len(start.Attributes) == 0 {
		return nil
	}

	for _, attr := range start.Attributes {
		item, found := n.attr(attr.Name.Value)
		if !found {
			//panic(fmt.Sprintf("不存在的属性 %s", attr.Name.Value))
			continue
		}

		if item.CanInterface() && item.Type().Implements(attrDecoderType) {
			if err := item.Interface().(AttrDecoder).DecodeXMLAttr(attr); err != nil {
				return err
			}
			continue
		} else if item.CanAddr() {
			pv := item.Addr()
			if pv.CanInterface() && pv.Type().Implements(attrDecoderType) {
				if err := pv.Interface().(AttrDecoder).DecodeXMLAttr(attr); err != nil {
					return err
				}
				continue
			}
		}

		if err := conv.Value(attr.Value.Value, item.Value); err != nil {
			return err
		}
	}

	return nil
}

func (n *node) decodeElements(p *Parser) error {
	for {
		t, err := p.Token()
		if err != nil {
			return err
		}
		if t == nil {
			return nil
		}

		switch elem := t.(type) {
		case *EndElement: // 找到当前对象的结束标签
			if elem.Name.Value != n.name {
				return p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
			}
			return nil
		case *CData:
			if !n.cdata.IsValid() {
				break
			}
			if err := conv.Value(elem.Value.Value, n.cdata.Value); err != nil {
				return err
			}
		case *String:
			if !n.content.IsValid() {
				break
			}
			if err := conv.Value(elem.Value, n.content.Value); err != nil {
				return err
			}
		case *StartElement:
			item, found := n.elem(elem.Name.Value)
			if !found { // 当子元素不存在于 elems 时，就有可能是当作内容处理的
				if !n.content.IsValid() {
					if err := findEndElement(p, elem); err != nil {
						return err
					}
					break
					// panic(fmt.Sprintf("不存在的子元素 %s", elem.Name.Value))
				}
				item = n.content
			}

			if err = n.decodeElement(p, elem, item); err != nil {
				return err
			}
		}
	} // end for
}

func (n *node) decodeElement(p *Parser, start *StartElement, v value) error {
	if v.CanInterface() && v.Type().Implements(decoderType) {
		return v.Interface().(Decoder).DecodeXML(p, start)
	} else if v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(decoderType) {
			return pv.Interface().(Decoder).DecodeXML(p, start)
		}
	}

	if isScalar(v.Value) {
		return decodeScalarElement(p, start, v.Value)
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Func, reflect.Chan, reflect.Array:
		panic(fmt.Sprintf("无效的 kind %s", v.Kind()))
	case reflect.Slice:
		return n.decodeSlice(p, start, v)
	default:
		oo := newNode(start.Name.Value, v.Value)
		return oo.decode(p, start)
	}
}

func (n *node) decodeSlice(p *Parser, start *StartElement, slice value) error {
	if start.Name.Value != slice.name {
		return findEndElement(p, start)
	}

	elem := reflect.New(slice.Type().Elem()).Elem()
	if elem.Kind() == reflect.Ptr {
		if elem.IsNil() {
			elem.Set(reflect.New(elem.Type().Elem()))
		}
	}

	if elem.CanInterface() && elem.Type().Implements(decoderType) {
		if err := elem.Interface().(Decoder).DecodeXML(p, start); err != nil {
			return err
		}
		goto RET
	} else if elem.CanAddr() {
		pv := elem.Addr()
		if pv.CanInterface() && pv.Type().Implements(decoderType) {
			if err := pv.Interface().(Decoder).DecodeXML(p, start); err != nil {
				return err
			}
			goto RET
		}
	}

	if isScalar(elem) {
		if err := decodeScalarElement(p, start, elem); err != nil {
			return err
		}
	} else {
		if err := newNode(start.Name.Value, elem).decode(p, start); err != nil {
			return err
		}
	}

RET:
	slice.Value.Set(reflect.Append(slice.Value, elem))
	return nil
}

func decodeScalarElement(p *Parser, start *StartElement, v reflect.Value) error {
	for {
		t, err := p.Token()
		if err != nil {
			return err
		}
		if t == nil {
			return nil
		}

		switch elem := t.(type) {
		case *EndElement:
			if elem.Name.Value != start.Name.Value {
				return p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
			}
			return nil
		case *String:
			if err = conv.Value(elem.Value, v); err != nil {
				return p.WithError(elem.Start, elem.End, err)
			}
		default:
			return p.NewError(start.Start, start.End, locale.ErrInvalidXML)
		}
	}
}

// 不相配，表示当前元素找不到与之相配的元素，需要忽略这个元素，
// 所以要过滤与 start 想匹配的结束符号才算结束。
func findEndElement(p *Parser, start *StartElement) error {
	level := 0
	for {
		t, err := p.Token()
		if err != nil {
			return err
		}
		if t == nil {
			return p.NewError(start.Start, start.End, locale.ErrInvalidXML) // 找不到相配的结束符号
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
