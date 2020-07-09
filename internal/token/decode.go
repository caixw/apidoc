// SPDX-License-Identifier: MIT

package token

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/node"
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
func Decode(h *core.MessageHandler, p *Parser, v interface{}, namespace string) {
	var hasRoot bool
	for {
		t, r, err := p.Token()
		if errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			h.Error(err)
			return
		}

		switch elem := t.(type) {
		case *StartElement:
			if hasRoot { // 多个根元素
				h.Error(p.NewError(elem.Start, elem.End, "", locale.ErrInvalidXML))
				return
			}
			hasRoot = true

			prefix := findPrefix(elem, namespace)
			vv := node.ParseValue(reflect.ValueOf(v))
			if err := decodeElement(p, elem, vv, prefix); err != nil {
				h.Error(err)
				return
			}
		case *Comment, *String, *Instruction: // 忽略注释和普通的文本内容
		default:
			h.Error(p.NewError(r.Start, r.End, "", locale.ErrInvalidXML))
			return
		}
	}
}

func findPrefix(start *StartElement, namespace string) string {
	for _, attr := range start.Attributes {
		if attr.Name.Prefix.Value == "xmlns" && attr.Value.Value == namespace {
			return attr.Name.Local.Value
		}
	}
	return ""
}

func decode(n *node.Node, p *Parser, start *StartElement, prefix string) (*EndElement, error) {
	if err := decodeAttributes(n, p, start, prefix); err != nil {
		return nil, err
	}

	if start.Close {
		return nil, decodeCheckOmitempty(n, p, start.Start, start.End)
	}

	end, err := decodeElements(n, p, prefix)
	if err != nil {
		return nil, err
	}
	return end, decodeCheckOmitempty(n, p, start.Start, end.End)
}

// 判断 omitempty 属性
func decodeCheckOmitempty(n *node.Node, p *Parser, start, end core.Position) error {
	for _, attr := range n.Attributes {
		if canNotEmpty(attr) {
			return p.NewError(start, end, attr.Name, locale.ErrRequired)
		}
	}
	for _, elem := range n.Elements {
		if canNotEmpty(elem) {
			return p.NewError(start, end, elem.Name, locale.ErrRequired)
		}
	}
	if n.CData != nil && canNotEmpty(n.CData) {
		return p.NewError(start, end, "cdata", locale.ErrRequired)
	}
	if n.Content != nil && canNotEmpty(n.Content) {
		return p.NewError(start, end, "content", locale.ErrRequired)
	}

	return nil
}

// 当前表示的值必须是一个非空值
func canNotEmpty(v *node.Value) bool {
	return v.Name != "" && // cdata 和 content 在未初始化时 name 字段为空值
		!v.Omitempty &&
		(!v.CanInterface() || is.Empty(v.Interface(), true))
}

// 将 start 的属性内容解码到 obj.Attributes 之中
func decodeAttributes(n *node.Node, p *Parser, start *StartElement, prefix string) error {
	if start == nil {
		return nil
	}

	for _, attr := range start.Attributes {
		item, found := n.Attribute(attr.Name.Local.Value)
		if !found || prefix != attr.Name.Prefix.Value { // 未找到或是命名空间不匹配
			continue
		}
		v := node.GetRealValue(item.Value)
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
			panic(fmt.Sprintf("当前属性 %s 未实现 AttrDecoder 接口", attr.Name))
		}

		if err := setAttributeValue(item.Value, item.Usage, p, attr); err != nil {
			return err
		}
	}

	return nil
}

func decodeElements(n *node.Node, p *Parser, prefix string) (*EndElement, error) {
	for {
		t, r, err := p.Token()
		if errors.Is(err, io.EOF) {
			// 应该只有 EndElement 才能返回，否则就不完整的 XML
			return nil, p.NewError(p.Position().Position, p.Position().Position, "", locale.ErrInvalidXML)
		} else if err != nil {
			return nil, err
		}

		switch elem := t.(type) {
		case *EndElement: // 找到当前对象的结束标签
			if (elem.Name.Local.Value == n.Value.Name) && (elem.Name.Prefix.Value == prefix) {
				return elem, nil
			}
			return nil, p.NewError(elem.Start, elem.End, n.Value.Name, locale.ErrNotFoundEndTag)
		case *CData:
			if n.CData != nil {
				setContentValue(n.CData.Value, reflect.ValueOf(elem))
			}
		case *String:
			if n.Content != nil {
				setContentValue(n.Content.Value, reflect.ValueOf(elem))
			}
		case *StartElement:
			item, found := n.Element(elem.Name.Local.Value)
			if !found || prefix != elem.Name.Prefix.Value {
				if err := findEndElement(p, elem); err != nil {
					return nil, err
				}
				break // 忽略不存在的子元素
			}
			if err = decodeElement(p, elem, item, prefix); err != nil {
				return nil, err
			}
		case *Comment: // 忽略注释内容
		default:
			return nil, p.NewError(r.Start, r.End, "", locale.ErrInvalidXML)
		}
	}
}

func setContentValue(target, source reflect.Value) {
	target = node.GetRealValue(target)
	source = node.GetRealValue(source)

	st := source.Type()
	num := st.NumField()
	for i := 0; i < num; i++ {
		name := st.Field(i).Name
		target.FieldByName(name).
			Set(source.Field(i))
	}
}

func decodeElement(p *Parser, start *StartElement, v *node.Value, prefix string) error {
	v.Value = node.GetRealValue(v.Value)
	k := v.Kind()
	switch {
	case k == reflect.Ptr, k == reflect.Func, k == reflect.Chan, k == reflect.Array, node.IsPrimitive(v.Value):
		panic(fmt.Sprintf("%s 是无效的类型", v.Value.Type()))
	case k == reflect.Slice:
		return decodeSlice(p, start, v, prefix)
	}

	end, impl, err := callDecodeXML(v.Value, p, start)
	if !impl {
		end, err = decode(node.New(start.Name.Local.Value, v.Value), p, start, prefix)
	}
	if err != nil {
		return err
	}
	return setTagValue(v.Value, v.Usage, p, start, end)
}

func decodeSlice(p *Parser, start *StartElement, slice *node.Value, prefix string) (err error) {
	// 不相配，表示当前元素找不到与之相配的元素，需要忽略这个元素，
	// 所以要过滤与 start 想匹配的结束符号才算结束。
	if !start.Close && (start.Name.Local.Value != slice.Name) {
		return findEndElement(p, start)
	}

	elem := reflect.New(slice.Type().Elem()).Elem()
	if elem.Kind() == reflect.Ptr && elem.IsNil() {
		elem.Set(reflect.New(elem.Type().Elem()))
	}

	end, impl, err := callDecodeXML(elem, p, start)
	if !impl {
		if node.IsPrimitive(elem) {
			panic(fmt.Sprintf("%s:%s 必须实现 Decoder 接口", slice.Name, elem.Type()))
		}
		end, err = decode(node.New(start.Name.Local.Value, elem), p, start, prefix)
	}
	if err != nil {
		return err
	}

	if err = setTagValue(elem, slice.Usage, p, start, end); err != nil {
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
		if errors.Is(err, io.EOF) {
			return p.NewError(start.Start, start.End, start.Name.String(), locale.ErrNotFoundEndTag) // 找不到相配的结束符号
		} else if err != nil {
			return err
		}

		switch elem := t.(type) {
		case *StartElement:
			if elem.Name.Equal(start.Name) {
				level++
			}
		case *EndElement:
			if level == 0 && start.Match(elem) {
				return nil
			}
			level--
		}
	}
}

func setTagValue(v reflect.Value, usage string, p *Parser, start *StartElement, end *EndElement) error {
	v = node.GetRealValue(v)
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("无效的 kind 类型: %s:%s", v.Type(), v.Kind()))
	}

	v.FieldByName(usageKeyName).Set(reflect.ValueOf(usage))
	v.FieldByName(elementTagName).Set(reflect.ValueOf(start.Name))
	if end == nil {
		v.FieldByName(rangeName).Set(reflect.ValueOf(start.Range))
	} else {
		v.FieldByName(rangeName).Set(reflect.ValueOf(core.Range{Start: start.Start, End: end.End}))
		v.FieldByName(elementTagEndName).Set(reflect.ValueOf(end.Name))
	}

	// Sanitize 在最后调用，可以保证 Sanitize 调用中可以取 v.Range 的值
	return callSanitizer(v, p)
}

func setAttributeValue(v reflect.Value, usage string, p *Parser, attr *Attribute) error {
	v = node.GetRealValue(v)
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("无效的 kind 类型: %s:%s", v.Type(), v.Kind()))
	}

	v.FieldByName(rangeName).Set(reflect.ValueOf(attr.Range))
	v.FieldByName(usageKeyName).Set(reflect.ValueOf(usage))
	v.FieldByName(attributeNameName).Set(reflect.ValueOf(attr.Name))

	return callSanitizer(v, p)
}

func callSanitizer(v reflect.Value, p *Parser) error {
	if v.CanInterface() && v.Type().Implements(sanitizerType) {
		return v.Interface().(Sanitizer).Sanitize(p)
	} else if v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(sanitizerType) {
			return pv.Interface().(Sanitizer).Sanitize(p)
		}
	}
	return nil
}
