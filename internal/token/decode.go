// SPDX-License-Identifier: MIT

package token

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/issue9/is"
	"golang.org/x/text/message"

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

type decoder struct {
	h      *core.MessageHandler
	p      *Parser
	prefix string // 命名空间所表示的前缀
}

func (d *decoder) message(t core.MessageType, start, end core.Position, field string, key message.Reference, v ...interface{}) bool {
	d.h.Message(t, d.p.NewError(start, end, field, key, v...))
	return false
}

func (d *decoder) error(err error) bool {
	d.h.Error(err)
	return false
}

// Decode 将 p 中的 XML 内容解码至 v 对象中
//
// namespace 如果不为空表示当前的 xml 所在的命名空间，只有该命名空间的元素才会被正确识别。
func Decode(h *core.MessageHandler, p *Parser, v interface{}, namespace string) {
	d := &decoder{
		h: h,
		p: p,
	}

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
				d.message(core.Erro, elem.Start, elem.End, "", locale.ErrInvalidXML)
				return
			}
			hasRoot = true

			d.prefix = findPrefix(elem, namespace)
			if !d.decodeElement(elem, node.ParseValue(reflect.ValueOf(v))) {
				return
			}
		case *Comment, *String, *Instruction: // 忽略注释和普通的文本内容
		default:
			d.message(core.Erro, r.Start, r.End, "", locale.ErrInvalidXML)
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

func (d *decoder) decode(n *node.Node, start *StartElement) *EndElement {
	d.decodeAttributes(n, start)

	if start.Close {
		d.checkOmitempty(n, start.Start, start.End)
		return nil
	}

	if end, ok := d.decodeElements(n); ok {
		d.checkOmitempty(n, start.Start, end.End)
		return end
	}
	return nil
}

// 判断 omitempty 属性
func (d *decoder) checkOmitempty(n *node.Node, start, end core.Position) {
	for _, attr := range n.Attributes {
		if canNotEmpty(attr) {
			d.message(core.Erro, start, end, attr.Name, locale.ErrRequired)
		}
	}
	for _, elem := range n.Elements {
		if canNotEmpty(elem) {
			d.message(core.Erro, start, end, elem.Name, locale.ErrRequired)
		}
	}
	if n.CData != nil && canNotEmpty(n.CData) {
		d.message(core.Erro, start, end, "cdata", locale.ErrRequired)
	}
	if n.Content != nil && canNotEmpty(n.Content) {
		d.message(core.Erro, start, end, "content", locale.ErrRequired)
	}
}

// 当前表示的值必须是一个非空值
func canNotEmpty(v *node.Value) bool {
	return v.Name != "" && // cdata 和 content 在未初始化时 name 字段为空值
		!v.Omitempty &&
		(!v.CanInterface() || is.Empty(v.Interface(), true))
}

// 将 start 的属性内容解码到 obj.Attributes 之中
func (d *decoder) decodeAttributes(n *node.Node, start *StartElement) {
	if start == nil {
		return
	}

	for _, attr := range start.Attributes {
		item, found := n.Attribute(attr.Name.Local.Value)
		if !found || d.prefix != attr.Name.Prefix.Value { // 未找到或是命名空间不匹配
			continue
		}
		v := node.GetRealValue(item.Value)
		v.Set(reflect.New(v.Type()).Elem())

		var impl bool
		if item.CanInterface() && item.Type().Implements(attrDecoderType) {
			if err := item.Interface().(AttrDecoder).DecodeXMLAttr(d.p, attr); err != nil {
				d.error(err)
			}
			impl = true
		} else if item.CanAddr() {
			pv := item.Addr()
			if pv.CanInterface() && pv.Type().Implements(attrDecoderType) {
				if err := pv.Interface().(AttrDecoder).DecodeXMLAttr(d.p, attr); err != nil {
					d.error(err)
				}
				impl = true
			}
		}

		if !impl {
			panic(fmt.Sprintf("当前属性 %s 未实现 AttrDecoder 接口", attr.Name))
		}

		d.setAttributeValue(item.Value, item.Usage, attr)
	}
}

func (d *decoder) decodeElements(n *node.Node) (end *EndElement, ok bool) {
	for {
		t, r, err := d.p.Token()
		if errors.Is(err, io.EOF) {
			// 应该只有 EndElement 才能返回，否则就不完整的 XML
			return nil, d.message(core.Erro, d.p.Position().Position, d.p.Position().Position, "", locale.ErrInvalidXML)
		} else if err != nil {
			return nil, d.error(err)
		}

		switch elem := t.(type) {
		case *EndElement: // 找到当前对象的结束标签
			if (elem.Name.Local.Value == n.Value.Name) && (elem.Name.Prefix.Value == d.prefix) {
				return elem, true
			}
			return nil, d.message(core.Erro, elem.Start, elem.End, n.Value.Name, locale.ErrNotFoundEndTag)
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
			if !found || d.prefix != elem.Name.Prefix.Value {
				if err := findEndElement(d.p, elem); err != nil {
					return nil, d.error(err)
				}

				d.message(core.Warn, elem.Name.Start, elem.Name.End, elem.Name.String(), locale.ErrInvalidTag)
				break // 忽略不存在的子元素
			}
			if !d.decodeElement(elem, item) {
				return nil, false
			}
		case *Comment: // 忽略注释内容
		default:
			return nil, d.message(core.Erro, r.Start, r.End, "", locale.ErrInvalidXML)
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

func (d *decoder) decodeElement(start *StartElement, v *node.Value) (ok bool) {
	v.Value = node.GetRealValue(v.Value)
	k := v.Kind()
	switch {
	case k == reflect.Ptr, k == reflect.Func, k == reflect.Chan, k == reflect.Array, node.IsPrimitive(v.Value):
		panic(fmt.Sprintf("%s 是无效的类型", v.Value.Type()))
	case k == reflect.Slice:
		return d.decodeSlice(start, v)
	}

	end, impl, err := callDecodeXML(v.Value, d.p, start)
	if !impl {
		end = d.decode(node.New(start.Name.Local.Value, v.Value), start)
	}
	if err != nil {
		return d.error(err)
	}

	d.setTagValue(v.Value, v.Usage, start, end)
	return true
}

func (d *decoder) decodeSlice(start *StartElement, slice *node.Value) (ok bool) {
	// 不相配，表示当前元素找不到与之相配的元素，需要忽略这个元素，
	// 所以要过滤与 start 想匹配的结束符号才算结束。
	if !start.Close && (start.Name.Local.Value != slice.Name) {
		if err := findEndElement(d.p, start); err != nil {
			return d.error(err)
		}
	}

	elem := reflect.New(slice.Type().Elem()).Elem()
	if elem.Kind() == reflect.Ptr && elem.IsNil() {
		elem.Set(reflect.New(elem.Type().Elem()))
	}

	end, impl, err := callDecodeXML(elem, d.p, start)
	if !impl {
		if node.IsPrimitive(elem) {
			panic(fmt.Sprintf("%s:%s 必须实现 Decoder 接口", slice.Name, elem.Type()))
		}
		end = d.decode(node.New(start.Name.Local.Value, elem), start)
	}
	if err != nil {
		return d.error(err)
	}

	d.setTagValue(elem, slice.Usage, start, end)
	slice.Value.Set(reflect.Append(slice.Value, elem))
	return true
}

// 调用 v 的 DecodeXML 接口方法
//
// 当 impl 为 true 时，err 表示的是 DecodeXML 接口返回的错误，否则 err 永远为 nil
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
	if start.Close {
		return nil
	}

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

func (d *decoder) setTagValue(v reflect.Value, usage string, start *StartElement, end *EndElement) {
	v = node.GetRealValue(v)
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("无效的 kind 类型: %s:%s", v.Type(), v.Kind()))
	}

	setFieldValue(v, usageKeyName, usage)
	setFieldValue(v, elementTagName, start.Name)
	if end == nil {
		setFieldValue(v, rangeName, start.Range)
	} else {
		setFieldValue(v, rangeName, core.Range{Start: start.Start, End: end.End})
		setFieldValue(v, elementTagEndName, end.Name)
	}

	// Sanitize 在最后调用，可以保证 Sanitize 调用中可以取 v.Range 的值
	if err := callSanitizer(v, d.p); err != nil {
		d.error(err)
	}
}

func (d *decoder) setAttributeValue(v reflect.Value, usage string, attr *Attribute) {
	v = node.GetRealValue(v)
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("无效的 kind 类型: %s:%s", v.Type(), v.Kind()))
	}

	setFieldValue(v, rangeName, attr.Range)
	setFieldValue(v, usageKeyName, usage)
	setFieldValue(v, attributeNameName, attr.Name)

	if err := callSanitizer(v, d.p); err != nil {
		d.error(err)
	}
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

func setFieldValue(v reflect.Value, name string, value interface{}) {
	vv := v.FieldByName(name)
	if vv.CanSet() {
		vv.Set(reflect.ValueOf(value))
	}
}
