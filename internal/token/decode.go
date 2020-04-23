// SPDX-License-Identifier: MIT

package token

import (
	"fmt"
	"reflect"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
)

// Decoder 实现从 p 中解码内容到当前对象
type Decoder interface {
	// 从 p 中读取内容并实例化到当前对象中
	//
	// 必须要同时从 p 中读取相应的 EndElement 才能返回。
	// end 表示 EndElement.End 的值。
	DecodeXML(p *Parser, start *StartElement) (end core.Position, err error)
}

// AttrDecoder 实现从 attr 中解码内容到当前对象
type AttrDecoder interface {
	DecodeXMLAttr(attr *Attribute) error
}

var (
	attrDecoderType = reflect.TypeOf((*AttrDecoder)(nil)).Elem()
	decoderType     = reflect.TypeOf((*Decoder)(nil)).Elem()
)

// Decode 将 p 中的 XML 内容解码至 v 对象中
func Decode(p *Parser, v interface{}) error {
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
			hasRoot = true

			o := newNode(elem.Name.Value, reflect.ValueOf(v))
			if _, err := o.decode(p, elem); err != nil {
				return err
			}
		case *EndElement:
			return p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
		}
	}
}

func (n *node) decode(p *Parser, start *StartElement) (core.Position, error) {
	if err := n.decodeAttributes(start); err != nil {
		return core.Position{}, err
	}

	if start.Close {
		return core.Position{}, nil
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

		var impl bool
		if item.CanInterface() && item.Type().Implements(attrDecoderType) {
			if err := item.Interface().(AttrDecoder).DecodeXMLAttr(attr); err != nil {
				return err
			}
			impl = true
		} else if item.CanAddr() {
			pv := item.Addr()
			if pv.CanInterface() && pv.Type().Implements(attrDecoderType) {
				if err := pv.Interface().(AttrDecoder).DecodeXMLAttr(attr); err != nil {
					return err
				}
				impl = true
			}
		}

		if !impl {
			panic(fmt.Sprintf("当前属性 %s 未实现 AttrDecoder 接口", attr.Name.Value))
		}

		initValue(item.Value, item.usage, attr.Start, attr.End)
	}

	return nil
}

func (n *node) decodeElements(p *Parser) (core.Position, error) {
	for {
		t, err := p.Token()
		if err != nil {
			return core.Position{}, err
		}
		if t == nil {
			return core.Position{}, nil
		}

		switch elem := t.(type) {
		case *EndElement: // 找到当前对象的结束标签
			if elem.Name.Value != n.name {
				return core.Position{}, p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
			}
			return elem.End, nil
		case *CData:
			if n.cdata.IsValid() {
				n.cdata.Set(getRealValue(reflect.ValueOf(elem)))
				initValue(n.cdata.Value, n.cdata.usage, elem.Start, elem.End)
			}
		case *String:
			if n.content.IsValid() {
				n.content.Set(getRealValue(reflect.ValueOf(elem)))
				initValue(n.content.Value, n.content.usage, elem.Start, elem.End)
			}
		case *StartElement:
			item, found := n.elem(elem.Name.Value)
			if !found {
				panic(fmt.Sprintf("不存在的子元素 %s", elem.Name.Value))
			}

			if err = decodeElement(p, elem, item); err != nil {
				return core.Position{}, err
			}
		}
	} // end for
}

func decodeElement(p *Parser, start *StartElement, v value) error {
	if v.CanInterface() && v.Type().Implements(decoderType) {
		end, err := v.Interface().(Decoder).DecodeXML(p, start)
		if err != nil {
			return err
		}
		initValue(v.Value, v.usage, start.Start, end)
		return nil
	} else if v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(decoderType) {
			end, err := pv.Interface().(Decoder).DecodeXML(p, start)
			if err != nil {
				return err
			}
			initValue(v.Value, v.usage, start.Start, end)
			return nil
		}
	}

	k := v.Kind()
	switch {
	case k == reflect.Ptr, k == reflect.Func, k == reflect.Chan, k == reflect.Array, isPrimitive(v.Value):
		panic(fmt.Sprintf("%s 是无效的类型", v.Value.Type().Name()))
	case k == reflect.Slice:
		return decodeSlice(p, start, v)
	default:
		end, err := newNode(start.Name.Value, v.Value).decode(p, start)
		if err != nil {
			return err
		}
		initValue(v.Value, v.usage, start.Start, end)
		return nil
	}
}

func decodeSlice(p *Parser, start *StartElement, slice value) (err error) {
	if start.Name.Value != slice.name {
		return findEndElement(p, start)
	}

	elem := reflect.New(slice.Type().Elem()).Elem()
	if elem.Kind() == reflect.Ptr {
		if elem.IsNil() {
			elem.Set(reflect.New(elem.Type().Elem()))
		}
	}

	var end core.Position
	if elem.CanInterface() && elem.Type().Implements(decoderType) {
		end, err = elem.Interface().(Decoder).DecodeXML(p, start)
		if err != nil {
			return err
		}
		goto RET
	} else if elem.CanAddr() {
		pv := elem.Addr()
		if pv.CanInterface() && pv.Type().Implements(decoderType) {
			end, err = pv.Interface().(Decoder).DecodeXML(p, start)
			if err != nil {
				return err
			}
			goto RET
		}
	}

	if isPrimitive(elem) {
		panic(fmt.Sprintf("%s:%s 必须实现 Decoder 接口", slice.name, elem.Type()))
	} else {
		if end, err = newNode(start.Name.Value, elem).decode(p, start); err != nil {
			return err
		}
	}

RET:
	initValue(elem, slice.usage, start.Start, end)
	slice.Value.Set(reflect.Append(slice.Value, elem))
	return nil
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

func initValue(v reflect.Value, usage string, start, end core.Position) {
	v = getRealValue(v)
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("无效的 kind 类型: %s:%s", v.Type(), v.Kind()))
	}

	v.FieldByName(rangeName).Set(reflect.ValueOf(core.Range{Start: start, End: end}))
	if usage != "" { // CDATA 和 content 节点类型的 usage 内容为空
		v.FieldByName(usageKeyName).Set(reflect.ValueOf(usage))
	}
}
